package model

import (
	"errors"

	"github.com/dapings/examples/go-programing-tour-2023/blog-service/pkg/app"
	"github.com/jinzhu/gorm"
)

type Article struct {
	*Model
	Title         string `json:"title"`
	Desc          string `json:"desc"`
	Content       string `json:"content"`
	CoverImageUrl string `json:"cover_image_url"`
	State         uint8  `json:"status"`
}

func (a Article) TableName() string {
	return "blog_article"
}

type ArticleSwagger struct {
	List  []*Article
	Pager *app.Pager
}

func (a Article) Create(db *gorm.DB) (*Article, error) {
	if err := db.Create(&a).Error; err != nil {
		return nil, err
	}

	return &a, nil
}

func (a Article) Update(db *gorm.DB, values interface{}) error {
	if err := db.Model(&a).Where("id = ? AND is_deleted = ?", a.ID, 0).Updates(values).Error; err != nil {
		return err
	}

	return nil
}

func (a Article) Get(db *gorm.DB) (Article, error) {
	var article Article
	db = db.Where("id = ? AND status = ? AND is_deleted = ?", a.ID, a.State, 0)
	err := db.First(&article).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return article, err
	}

	return article, nil
}

func (a Article) Delete(db *gorm.DB) error {
	if err := db.Where("id = ? AND is_deleted = ?", a.Model.ID, 0).Delete(&a).Error; err != nil {
		return err
	}

	return nil
}

func (a Article) CountByTagID(db *gorm.DB, tagID uint32) (int, error) {
	var count int
	err := db.Table(ArticleTag{}.TableName()+" AS at").
		Joins("LEFT JOIN `"+Tag{}.TableName()+"` AS t ON at.tag_id = t.id").
		Joins("LEFT JOIN `"+Article{}.TableName()+"` AS ar ON at.article_id = ar.id").
		Where("at.`tag_id` = ? AND ar.status = ? AND ar.is_deleted = ?", tagID, a.State, 0).
		Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}
