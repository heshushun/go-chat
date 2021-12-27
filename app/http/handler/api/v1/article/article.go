package article

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-chat/app/http/request"
	"go-chat/app/http/response"
	"go-chat/app/pkg/auth"
	"go-chat/app/pkg/filesystem"
	"go-chat/app/pkg/slice"
	"go-chat/app/pkg/strutil"
	"go-chat/app/pkg/timeutil"
	"go-chat/app/pkg/utils"
	"go-chat/app/service/note"
)

type Article struct {
	service             *note.ArticleService
	fileSystem          *filesystem.Filesystem
	articleAnnexService *note.ArticleAnnexService
}

func NewArticleHandler(service *note.ArticleService, fileSystem *filesystem.Filesystem, articleAnnexService *note.ArticleAnnexService) *Article {
	return &Article{service, fileSystem, articleAnnexService}
}

// List 文章列表
func (c *Article) List(ctx *gin.Context) {
	params := &request.ArticleListRequest{}

	if err := ctx.ShouldBindQuery(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	items, err := c.service.List(ctx.Request.Context(), auth.GetAuthUserID(ctx), params)
	if err != nil {
		response.BusinessError(ctx, err)
		return
	}

	list := make([]map[string]interface{}, 0)
	for _, item := range items {
		list = append(list, map[string]interface{}{
			"abstract":    item.Abstract,
			"class_id":    item.ClassId,
			"created_at":  timeutil.FormatDatetime(item.CreatedAt),
			"id":          item.Id,
			"image":       item.Image,
			"is_asterisk": item.IsAsterisk,
			"status":      item.Status,
			"tags_id":     item.TagsId,
			"title":       item.Title,
			"updated_at":  timeutil.FormatDatetime(item.UpdatedAt),
		})
	}

	response.SuccessPaginate(ctx, list, 1, 1000, len(items))
}

// Detail 文章详情
func (c *Article) Detail(ctx *gin.Context) {
	params := &request.ArticleDetailRequest{}

	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	uid := auth.GetAuthUserID(ctx)

	detail, err := c.service.Detail(ctx.Request.Context(), uid, params.ArticleId)
	if err != nil {
		response.BusinessError(ctx, "笔记不存在")
		return
	}

	tags := make([]map[string]interface{}, 0)

	for _, tagId := range slice.ParseIds(detail.TagsId) {
		tags = append(tags, map[string]interface{}{
			"id": tagId,
		})
	}

	files := make([]map[string]interface{}, 0)
	items, err := c.articleAnnexService.AnnexList(ctx, uid, params.ArticleId)
	if err == nil {
		for _, item := range items {
			files = append(files, map[string]interface{}{
				"id":            item.Id,
				"file_suffix":   item.FileSuffix,
				"file_size":     item.FileSize,
				"original_name": item.OriginalName,
				"created_at":    timeutil.FormatDatetime(item.CreatedAt),
			})
		}
	}

	response.Success(ctx, gin.H{
		"id":          detail.Id,
		"class_id":    detail.ClassId,
		"title":       detail.Title,
		"md_content":  detail.MdContent,
		"content":     detail.Content,
		"is_asterisk": detail.IsAsterisk,
		"status":      detail.Status,
		"created_at":  timeutil.FormatDatetime(detail.CreatedAt),
		"updated_at":  timeutil.FormatDatetime(detail.UpdatedAt),
		"tags":        tags,
		"files":       files,
	})
}

// Class 添加或编辑文章
func (c *Article) Edit(ctx *gin.Context) {
	var (
		err    error
		params = &request.ArticleEditRequest{}
		uid    = auth.GetAuthUserID(ctx)
	)

	if err = ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	if params.ArticleId == 0 {
		params.ArticleId, err = c.service.Create(ctx.Request.Context(), uid, params)
	} else {
		err = c.service.Update(ctx.Request.Context(), uid, params)
	}

	if err != nil {
		response.BusinessError(ctx, err)
		return
	}

	response.Success(ctx, gin.H{
		"id": params.ArticleId,
	})
}

// Delete 删除文章
func (c *Article) Delete(ctx *gin.Context) {

	params := &request.ArticleDeleteRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	err := c.service.UpdateStatus(ctx.Request.Context(), auth.GetAuthUserID(ctx), params.ArticleId, 2)
	if err != nil {
		response.BusinessError(ctx, err)
		return
	}

	response.Success(ctx, gin.H{})
}

// Recover 恢复文章
func (c *Article) Recover(ctx *gin.Context) {
	params := &request.ArticleRecoverRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	err := c.service.UpdateStatus(ctx.Request.Context(), auth.GetAuthUserID(ctx), params.ArticleId, 1)
	if err != nil {
		response.BusinessError(ctx, err)
		return
	}

	response.Success(ctx, gin.H{})
}

// Upload 文章图片上传
func (c *Article) Upload(ctx *gin.Context) {
	file, err := ctx.FormFile("image")
	if err != nil {
		response.InvalidParams(ctx, "image 字段必传！")
		return
	}

	if !slice.InStr(strutil.FileSuffix(file.Filename), []string{"png", "jpg", "jpeg", "gif", "webp"}) {
		response.InvalidParams(ctx, "上传文件格式不正确,仅支持 png、jpg、jpeg、gif 和 webp")
		return
	}

	// 判断上传文件大小（5M）
	if file.Size > 5<<20 {
		response.InvalidParams(ctx, "上传文件大小不能超过5M！")
		return
	}

	stream, err := filesystem.ReadMultipartStream(file)
	if err != nil {
		response.BusinessError(ctx, "文件上传失败")
		return
	}

	ext := strutil.FileSuffix(file.Filename)
	m := utils.ReadFileImage(bytes.NewReader(stream))

	filePath := fmt.Sprintf("public/media/image/note/%s/%s", timeutil.DateNumber(), strutil.GenImageName(ext, m.Width, m.Height))

	if err := c.fileSystem.Default.Write(stream, filePath); err != nil {
		response.BusinessError(ctx, "文件上传失败")
		return
	}

	response.Success(ctx, gin.H{"save_path": c.fileSystem.Default.PublicUrl(filePath)})
}

// Move 文章移动
func (c *Article) Move(ctx *gin.Context) {

	params := &request.ArticleMoveRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	err := c.service.Move(ctx.Request.Context(), auth.GetAuthUserID(ctx), params)
	if err != nil {
		response.BusinessError(ctx, err)
		return
	}

	response.Success(ctx, gin.H{})
}

// Asterisk 标记文章
func (c Article) Asterisk(ctx *gin.Context) {

	params := &request.ArticleAsteriskRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	err := c.service.Asterisk(ctx.Request.Context(), auth.GetAuthUserID(ctx), params)
	if err != nil {
		response.BusinessError(ctx, err)
		return
	}

	response.Success(ctx, gin.H{})
}

// UpdateTag 文章标签
func (c *Article) UpdateTag(ctx *gin.Context) {
	params := &request.ArticleTagsRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	err := c.service.UpdateTag(ctx.Request.Context(), auth.GetAuthUserID(ctx), params)
	if err != nil {
		response.BusinessError(ctx, err)
		return
	}

	response.Success(ctx, gin.H{})
}

// ForeverDelete 永久删除文章
func (c *Article) ForeverDelete(ctx *gin.Context) {
	params := &request.ArticleForeverDeleteRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	err := c.service.ForeverDelete(ctx.Request.Context(), auth.GetAuthUserID(ctx), params)
	if err != nil {
		response.BusinessError(ctx, err)
		return
	}

	response.Success(ctx, gin.H{})
}
