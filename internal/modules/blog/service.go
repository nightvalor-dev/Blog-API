package blog

import (
	"Project2-v7/internal/modules/category"
	"Project2-v7/internal/modules/media"
	"Project2-v7/internal/modules/tag"
	"Project2-v7/internal/shared/redis"
	"Project2-v7/pkg/utils"
	"context"
)

type BlogService struct {
	blogrepo     BlogRepository
	categoryRepo category.CategoryRepository
	tagRepo      tag.TagRepository
	mediaRepo    media.MediaRepository
	cache        *redis.Cache
}

func NewBlogService(
	blogrepo BlogRepository,
	categoryRepo category.CategoryRepository,
	tagRepo tag.TagRepository,
	mediaRepo media.MediaRepository,
	cache *redis.Cache) *BlogService {

	return &BlogService{
		blogrepo:     blogrepo,
		categoryRepo: categoryRepo,
		tagRepo:      tagRepo,
		mediaRepo:    mediaRepo,
		cache:        cache,
	}
}

func (s *BlogService) Create(ctx context.Context, req CreateBlogRequest) error {
	blogId, err := s.blogrepo.Create(req)
	if err != nil {
		return err
	}

	for _, c := range req.Category {
		categoryId, err := s.categoryRepo.FindOrCreate(category.CategoryRequest{
			CategoryName: c.CategoryName,
			Description:  c.Description,
		})
		if err != nil {
			return err
		}

		if err = s.blogrepo.AddCategory(blogId, categoryId); err != nil {
			return err
		}
	}

	for _, t := range req.Tags {
		tagId, err := s.tagRepo.FindOrCreate(tag.TagRequest{TagName: t.TagName})

		if err != nil {
			return err
		}

		if err = s.blogrepo.AddTag(blogId, tagId); err != nil {
			return err
		}
	}

	_ = s.cache.DeleteByPrefix(ctx, redis.BlogList) // new post changes page composition
	return nil
}

func (s *BlogService) GetAllFiltered(ctx context.Context, filter BlogFilter, pagination utils.PaginationParams) (BlogListResponse, error) {
	blogs, total, err := s.blogrepo.GetAllFiltered(filter, pagination)
	if err != nil {
		return BlogListResponse{}, err
	}

	data, err := s.buildResponseList(blogs)
	if err != nil {
		return BlogListResponse{}, err
	}

	return BlogListResponse{
		Data:       data,
		TotalCount: total,
		Page:       pagination.Page,
		Limit:      pagination.Limit,
	}, nil
}

func (s *BlogService) GetAllCached(ctx context.Context, pagination utils.PaginationParams) (BlogListResponse, error) {
	key := redis.BlogListKey(pagination.Page, pagination.Limit)

	var result BlogListResponse
	if err := s.cache.Get(ctx, key, &result); err == nil {
		return result, nil
	}

	defaultFilter := BlogFilter{}
	blogs, total, err := s.blogrepo.GetAllFiltered(defaultFilter, pagination)
	if err != nil {
		return BlogListResponse{}, err
	}

	data, err := s.buildResponseList(blogs)
	if err != nil {
		return BlogListResponse{}, err
	}

	result = BlogListResponse{
		Data:       data,
		TotalCount: total,
		Page:       pagination.Page,
		Limit:      pagination.Limit,
	}

	_ = s.cache.Set(ctx, key, result, redis.DefaultTTL)

	return result, nil
}

func (s *BlogService) GetById(ctx context.Context, id int) (BlogResponse, error) {
	key := redis.BlogKey(id)
	var response BlogResponse

	if err := s.cache.Get(ctx, key, &response); err == nil {
		return response, nil
	}

	b, err := s.blogrepo.GetById(id)
	if err != nil {
		return BlogResponse{}, err
	}

	response, err = s.buildResponse(b)
	if err != nil {
		return BlogResponse{}, err
	}

	_ = s.cache.Set(ctx, key, response, redis.DefaultTTL)

	return response, nil
}

func (s *BlogService) Update(ctx context.Context, id int, req UpdateBlogRequest) error {
	err := s.blogrepo.Update(id, req)
	if err != nil {
		return err
	}

	_ = s.cache.Delete(ctx, redis.BlogKey(id))
	_ = s.cache.DeleteByPrefix(ctx, redis.BlogList)
	return nil
}

func (s *BlogService) Delete(ctx context.Context, id int) error {
	if err := s.blogrepo.Delete(id); err != nil {
		return err
	}
	_ = s.cache.Delete(ctx, redis.BlogKey(id))
	_ = s.cache.DeleteByPrefix(ctx, redis.BlogList)
	return nil
}

func (s *BlogService) buildResponseList(blogs []Blog) ([]BlogResponse, error) {
	data := make([]BlogResponse, 0, len(blogs))
	for _, bl := range blogs {
		response, err := s.buildResponse(bl)
		if err != nil {
			return nil, err
		}
		data = append(data, response)
	}
	return data, nil
}

func (s *BlogService) buildResponse(b Blog) (BlogResponse, error) {
	cats, err := s.categoryRepo.GetByBlogId(b.BlogId)
	if err != nil {
		return BlogResponse{}, err
	}

	tags, err := s.tagRepo.GetByBlogId(b.BlogId)
	if err != nil {
		return BlogResponse{}, err
	}

	mediaList, err := s.mediaRepo.GetByBlogId(b.BlogId)
	if err != nil {
		return BlogResponse{}, err
	}

	catDTOs := make([]CategoryDTOLocal, len(cats))
	for i, c := range cats {
		catDTOs[i] = CategoryDTOLocal{
			CategoryName: c.CategoryName,
			Description:  c.Description,
		}
	}

	tagDTOs := make([]TagDTOLocal, len(tags))
	for i, t := range tags {
		tagDTOs[i] = TagDTOLocal{TagName: t.TagName}
	}

	mediaDTOs := make([]MediaDTOLocal, len(mediaList))
	for i, m := range mediaList {
		mediaDTOs[i] = MediaDTOLocal{Url: m.Url}
	}

	return BlogResponse{
		Title:    b.Title,
		Content:  b.Content,
		Status:   b.Status,
		Category: catDTOs,
		Tags:     tagDTOs,
		Images:   mediaDTOs,
	}, nil
}
