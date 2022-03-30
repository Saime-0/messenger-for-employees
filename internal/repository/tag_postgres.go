package repository

import (
	"database/sql"
	"github.com/saime-0/http-cute-chat/graph/model"
)

type TagsRepo struct {
	db *sql.DB
}

func NewTagsRepo(db *sql.DB) *TagsRepo {
	return &TagsRepo{
		db: db,
	}
}

func (r *TagsRepo) Tags(params *model.Params) (*model.Tags, error) {

	tags := &model.Tags{
		Tags: []*model.Tag{},
	}

	var rows, err = r.db.Query(`
		SELECT t.tag_id, t.name
		FROM tags t
		LIMIT $1
		OFFSET $2
		`,
		params.Limit,
		params.Offset,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		m := new(model.Tag)
		if err = rows.Scan(&m.TagID, &m.Name); err != nil {
			return nil, err
		}
		tags.Tags = append(tags.Tags, m)
	}

	return tags, nil
}
