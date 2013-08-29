package meta

import (
	"time"
)

// 存储每个目录、标签
type Terms struct {
	Term_id   uint64 //分类id  //BIGINT( 20 ) UNSIGNED NOT NULL AUTO_INCREMENT ,
	User_id   uint64 //User_id //BIGINT( 20 ) UNSIGNED NOT NULL,
	Term_name string //分类名  //VARCHAR( 200 ) NOT NULL DEFAULT '',
	Term_slug string //缩略名  //VARCHAR( 200 ) NOT NULL DEFAULT '',
}

// 存储每个目录、标签所对应的分类
// 分类信息,是对 Terms 中的信息的关系信息补充，
// 有所属类型(category,link_category,tag)，详细描述所拥有文章(链接)数量。
type Termtaxonomy struct {
	Termtaxonomy_id uint64 //分类方法id                            //BIGINT( 20 ) UNSIGNED NOT NULL AUTO_INCREMENT ,
	Term_id         uint64 //Term_id                               //BIGINT( 20 ) UNSIGNED NOT NULL DEFAULT 0,
	Taxonomy        string //taxonomy:分类方法(category/post_tag)  //VARCHAR( 32 ) NOT NULL DEFAULT '',
	Description     string //说明                                  //LONGTEXT NOT NULL ,
	Parent          uint64 //所属父分类方法id                      //BIGINT( 20 ) UNSIGNED NOT NULL DEFAULT 0,
	Count           uint64 //文章数统计                            //BIGINT( 20 ) UNSIGNED NOT NULL DEFAULT 0,
}

// 存储每个文章、链接和对应分类的关系
// 把 posts和links这些对象和term_taxonomy表中的term_taxonomy_id联系起来的关系表，
// object_id是与不同的对象关联，例如posts中的id（links中的link_id）等，
// termtaxonomy_id就是关联 termtaxonomy中的termtaxonomy_id
type Termrelationships struct {
	Object_id       uint64 //对应文章id/链接id //BIGINT( 20 ) UNSIGNED NOT NULL DEFAULT 0,
	Termtaxonomy_id uint64 //对应分类方法id    //BIGINT( 20 ) UNSIGNED NOT NULL DEFAULT 0,
	Term_order      uint32 //排序              //INT( 11 ) UNSIGNED NOT NULL DEFAULT 0,
}

// 评论中的额外数据
type Commentmeta struct {
	Commentmeta_id uint64 //BIGINT( 20 ) UNSIGNED NOT NULL AUTO_INCREMENT ,
	Comment_id     uint64 //BIGINT( 20 ) UNSIGNED NOT NULL,
	User_id        uint64 //冗余记录 post User_id //BIGINT( 20 ) UNSIGNED NOT NULL,
	Meta_key       string //VARCHAR( 255 ) DEFAULT NULL ,
	Meta_value     string //LONGTEXT,
}

// 评论
type Comments struct {
	Comment_id      uint64    //自增唯一id                //BIGINT( 20 ) UNSIGNED NOT NULL AUTO_INCREMENT ,
	Post_id         uint64    //对应文章id                //BIGINT( 20 ) UNSIGNED NOT NULL DEFAULT '0',
	Comment_date    time.Time //评论时间                  //TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	Comment_ip      string    //评论者ip                  //VARCHAR( 100 ) NOT NULL DEFAULT '',
	Comment_content string    //评论正文                  //TEXT NOT NULL ,
	Comment_vetoed  string    //评论被否决的理由          //VARCHAR( 20 ) NOT NULL DEFAULT '',
	Comment_type    string    //评论类型(pingback/普通)   //VARCHAR( 20 ) NOT NULL DEFAULT '',
	Comment_parent  uint64    //父评论id引用              //BIGINT( 20 ) UNSIGNED NOT NULL DEFAULT '0',
	User_nicename   string    //评论者                    //VARCHAR( 50 ) NOT NULL DEFAULT '',
	User_id         uint64    //评论者用户id（不一定存在）//BIGINT( 20 ) UNSIGNED NOT NULL DEFAULT '0',
}

// 存储友情链接（Blogroll）
type Links struct {
	Link_id          uint64    //自增唯一id     //BIGINT( 20 ) UNSIGNED NOT NULL AUTO_INCREMENT ,
	User_id          uint64    //用户id         //BIGINT( 20 ) UNSIGNED NOT,
	Link_url         string    //链接url        //VARCHAR( 255 ) NOT NULL,
	Link_rss         string    //链接rss地址    //VARCHAR( 255 ) NOT NULL DEFAULT '',
	Link_title       string    //链接标题       //VARCHAR( 255 ) NOT NULL DEFAULT '',
	Link_image       string    //链接图片       //VARCHAR( 255 ) NOT NULL DEFAULT '',
	Link_target      string    //链接打开方式   //VARCHAR( 25 ) NOT NULL DEFAULT '',
	Link_description string    //链接描述       //VARCHAR( 255 ) NOT NULL DEFAULT '',
	Link_rel         string    //xfn关系        //VARCHAR( 255 ) NOT NULL DEFAULT '',
	Link_notes       string    //xfn注释        //MEDIUMTEXT NOT NULL ,
	Link_updated     time.Time //更新时间       //TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	Link_rating      uint32    //评分等级       //INT( 11 ) UNSIGNED NOT NULL DEFAULT '0',
	Link_vetoed      string    //是否不可见     //VARCHAR( 20 ) NOT NULL DEFAULT '',
}

// 发布的元数据（包括页面、上传文件、修订）的元数据
type Postmeta struct {
	Postmeta_id uint64 //自增唯一ID //BIGINT( 20 ) UNSIGNED NOT NULL AUTO_INCREMENT ,
	Post_id     uint64 //对应文章ID //BIGINT( 20 ) UNSIGNED NOT NULL DEFAULT '0',
	User_id     uint64 //对应作者id //BIGINT( 20 ) UNSIGNED NOT NULL DEFAULT '0',
	Meta_key    string //键名       //VARCHAR( 255 ) DEFAULT NULL ,
	Meta_value  string //键值       //LONGTEXT,
}

// 存储文章（包括页面、上传文件、修订）
type Posts struct {
	Post_id        uint64    //自增唯一id                              //BIGINT( 20 ) UNSIGNED NOT NULL AUTO_INCREMENT ,
	User_id        uint64    //对应作者id                              //BIGINT( 20 ) UNSIGNED NOT NULL DEFAULT '0',
	Post_title     string    //标题                                    //VARCHAR( 200 ) NOT NULL DEFAULT '',
	Post_slug      string    //文章缩略名                              //VARCHAR( 200 ) NOT NULL DEFAULT '',
	Post_content   string    //正文                                    //LONGTEXT NOT NULL ,
	Post_excerpt   string    //摘录                                    //TEXT NOT NULL ,
	Post_status    string    //文章状态（publish/auto-draft/inherit等）//VARCHAR( 20 ) NOT NULL DEFAULT 'publish',
	Comment_vetoed string    //拒绝评论                                //VARCHAR( 20 ) NOT NULL DEFAULT '',
	Ping_vetoed    string    //拒绝ping                                //VARCHAR( 20 ) NOT NULL DEFAULT '',
	Post_password  string    //文章密码                                //VARCHAR( 32 ) NOT NULL DEFAULT '',
	To_ping        string    //Pingback                                //TEXT NOT NULL ,
	Pinged         string    //已经ping过的链接                        //TEXT NOT NULL ,
	Post_date      time.Time //发布时间                                //TIMESTAMP NOT NULL DEFAULT '0000-00-00 00:00:00',
	Post_modified  time.Time //修改时间                                //TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	Post_parent    uint64    //父文章，主要用于分页                    //BIGINT( 20 ) UNSIGNED NOT NULL DEFAULT '0',
	Guid           string    //未知(url)                               //VARCHAR( 255 ) NOT NULL DEFAULT '',
	Menu_order     uint32    //排序id                                  //INT( 11 ) UNSIGNED NOT NULL DEFAULT '0',
	Post_mime_type string    //mime类型                                //VARCHAR( 100 ) NOT NULL DEFAULT '',
	Comment_count  uint64    //评论总数                                //BIGINT( 20 ) UNSIGNED NOT NULL DEFAULT '0',
}

type UserStatus int32

const (
	UserStatusOK      UserStatus = iota //唯一正常状态
	UserStatusHIDE                      //用户:自我隐蔽,可以登录,数据保留,关闭一切数据展现
	UserStatusFREE                      //用户:自我解除,不能登录,数据保留,关闭一切数据展现
	UserStatusEND                       //用户:自我解除,不能登录,请求删除所有数据
	UserStatusSPAM                      //管理:鉴定为spam,可以登录,数据保留,关闭一切数据展现
	UserStatusDISPUTE                   //管理:鉴定为dispute,可以登录,数据保留,关闭一切数据展现
	UserStatusSHUT                      //管理:鉴定为shut,可以登录,即将被彻底关闭,数据保留,关闭一切数据展现
	UserStatusDOWN                      //管理:鉴定为down,彻底关闭,数据保留,关闭一切数据展现
	UserStatusEof                       //这是一个功能常量为了判断最大值
)

// 用户信息，用户必须设定博客子域名
type Users struct {
	User_id       uint64    //自增唯一id //BIGINT( 20 ) UNSIGNED NOT NULL AUTO_INCREMENT ,
	User_login    string    //md5 登录名 //CHAR( 32 ) NOT NULL,
	User_pass     string    //md5 密码   //CHAR( 32 ) NOT NULL,
	User_nicename string    //显示名称   //VARCHAR( 50 ) NOT NULL,
	User_date     time.Time //注册时间   //DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	User_status   uint32    //用户状态   //INT( 11 ) UNSIGNED NOT NULL DEFAULT '0',
	Site          string    //博客子域名 //VARCHAR( 20 ) NOT NULL DEFAULT '',
}

// 用户额外数据
type Usermeta struct {
	Usermeta_id uint64 //自增唯一id //BIGINT( 20 ) UNSIGNED NOT NULL AUTO_INCREMENT ,
	User_id     uint64 //对应用户id //BIGINT( 20 ) UNSIGNED NOT NULL DEFAULT '0',
	Meta_tag    string //分类       //VARCHAR( 20 ) DEFAULT NULL ,
	Meta_key    string //键名       //VARCHAR( 200 ) DEFAULT NULL ,
	Meta_value  string //键值       //LONGTEXT,
}

// 组织成员
type Members struct {
	User_id       uint64 //对应users的user_id  //BIGINT( 20 ) UNSIGNED NOT NULL,
	Member_id     uint64 //成员的user_id       //BIGINT( 20 ) UNSIGNED NOT NULL,
	Member_status uint32 //状态 //INT( 11 ) UNSIGNED NOT NULL DEFAULT '0',
}

// 站点额外数据
type Sitemeta struct {
	Sitemeta_id uint64 //BIGINT( 20 ) UNSIGNED NOT NULL AUTO_INCREMENT ,
	User_id     uint64 //BIGINT( 20 ) UNSIGNED NOT NULL,
	Meta_key    string //VARCHAR( 255 ) DEFAULT NULL ,
	Meta_value  string //LONGTEXT,
}
