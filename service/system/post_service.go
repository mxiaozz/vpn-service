package system

import (
	"github.com/pkg/errors"
	"vpn-web.funcworks.net/gb"
	"vpn-web.funcworks.net/model"
	"vpn-web.funcworks.net/model/entity"
	"xorm.io/builder"
	"xorm.io/xorm"
)

var PostService = &postService{}

type postService struct {
}

func (ps *postService) GetALlPosts() ([]entity.SysPost, error) {
	var posts = []entity.SysPost{}
	err := gb.DB.Find(&posts)
	return posts, err
}

// 岗位管理，获取岗位列表（包括岗位查询/分页）
func (ps *postService) GetPostListPage(post *entity.SysPost, page *model.Page[entity.SysPost]) error {
	return gb.SelectPage(page, func(sql *builder.Builder) builder.Cond {
		sql.Select("*").From("sys_post").
			Where(builder.If(post.PostCode != "", builder.Like{"post_code", post.PostCode}).
				And(builder.If(post.PostName != "", builder.Like{"post_name", post.PostName})).
				And(builder.If(post.Status != "", builder.Eq{"status": post.Status})))
		return builder.Expr("post_sort asc")
	})
}

func (ps *postService) GetUserPostList(userId int64) ([]entity.SysPost, error) {
	var posts = []entity.SysPost{}
	err := gb.DB.Table("sys_post").Alias("p").
		Join("left", []string{"sys_user_post", "up"}, "up.post_id = p.post_id").
		Join("left", []string{"sys_user", "u"}, "u.user_id = up.user_id").
		Where("u.user_id = ?", userId).
		Find(&posts)
	return posts, err
}

func (ps *postService) GetPost(postId int64) (*entity.SysPost, error) {
	var post entity.SysPost
	if exist, err := gb.DB.Where("post_id = ?", postId).Get(&post); err != nil {
		return nil, err
	} else if !exist {
		return nil, nil
	}
	return &post, nil
}

func (ps *postService) AddPost(post *entity.SysPost) error {
	_, err := gb.DB.Insert(post)
	return err
}

func (ps *postService) UpdatePost(post *entity.SysPost) error {
	_, err := gb.DB.Where("post_id = ?", post.PostId).Update(post)
	return err
}

func (ps *postService) DeletePosts(postIds []int64) error {
	return gb.Tx(func(dbSession *xorm.Session) error {
		for _, postId := range postIds {
			if exist, err := ps.checkPostRsUser(postId); err != nil {
				return err
			} else if exist {
				return errors.New("岗位存在用户关联，不能删除")
			}
			if _, err := gb.DB.Table("sys_post").Where("post_id = ?", postId).Delete(); err != nil {
				return err
			}
		}
		return nil
	})
}

func (ps *postService) checkPostRsUser(postId int64) (bool, error) {
	return gb.DB.Table("sys_user_post").Where("post_id = ?", postId).Exist()
}
