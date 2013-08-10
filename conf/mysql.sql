
CREATE TABLE IF NOT EXISTS `terms` (
  `term_id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `term_name` varchar(200) NOT NULL DEFAULT '',
  `term_slug` varchar(200) NOT NULL DEFAULT '',
  `user_id` bigint(20) unsigned NOT NULL,
  PRIMARY KEY (`term_id`),
  UNIQUE KEY `term_slug` (`user_id`,`term_slug`),
  KEY `term_name` (`term_name`)
) DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `termtaxonomy` (
  `termtaxonomy_id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `term_id` bigint(20) unsigned NOT NULL DEFAULT '0',
  `taxonomy` varchar(32) NOT NULL DEFAULT '',
  `description` longtext NOT NULL,
  `parent` bigint(20) unsigned NOT NULL DEFAULT '0',
  `count` bigint(20) unsigned NOT NULL DEFAULT '0',
  `blog_id` bigint(20) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`termtaxonomy_id`),
  UNIQUE KEY `term_id_taxonomy` (`term_id`,`taxonomy`),
  KEY `taxonomy` (`taxonomy`)
) DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `termrelationships` (
  `object_id` bigint(20) unsigned NOT NULL DEFAULT '0',
  `termtaxonomy_id` bigint(20) unsigned NOT NULL DEFAULT '0',
  `term_order` int(11) unsigned NOT NULL DEFAULT '0',
  `blog_id` bigint(20) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`object_id`,`termtaxonomy_id`),
  KEY `termtaxonomy_id` (`termtaxonomy_id`)
) DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `commentmeta` (
  `commentmeta_id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `comment_id` bigint(20) unsigned NOT NULL DEFAULT '0',
  `meta_key` varchar(255) DEFAULT NULL,
  `meta_value` longtext,
  `blog_id` bigint(20) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`commentmeta_id`),
  KEY `comment_id` (`comment_id`),
  KEY `meta_key` (`meta_key`)
) DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `comments` (
  `comment_id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `comment_post_id` bigint(20) unsigned NOT NULL DEFAULT '0',
  `comment_author` varchar(50) NOT NULL DEFAULT '',
  `comment_author_email` varchar(100) NOT NULL DEFAULT '',
  `comment_author_url` varchar(200) NOT NULL DEFAULT '',
  `comment_author_ip` varchar(100) NOT NULL DEFAULT '',
  `comment_date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `comment_content` text NOT NULL,
  `comment_vetoed` varchar(20) NOT NULL DEFAULT '',
  `comment_agent` varchar(255) NOT NULL DEFAULT '',
  `comment_type` varchar(20) NOT NULL DEFAULT '',
  `comment_parent` bigint(20) unsigned NOT NULL DEFAULT '0',
  `user_id` bigint(20) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`comment_id`),
  KEY `comment_post_id` (`comment_post_id`),
  KEY `comment_approved_date_gmt` (`comment_vetoed`),
  KEY `comment_date` (`comment_date`),
  KEY `comment_parent` (`comment_parent`)
) DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `links` (
  `link_id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) unsigned NOT NULL,
  `link_url` varchar(255) NOT NULL,
  `link_title` varchar(255) NOT NULL DEFAULT '',
  `link_image` varchar(255) NOT NULL DEFAULT '',
  `link_target` varchar(25) NOT NULL DEFAULT '',
  `link_description` varchar(255) NOT NULL DEFAULT '',
  `link_vetoed` varchar(20) NOT NULL DEFAULT '',
  `link_rating` int(11) unsigned NOT NULL DEFAULT '0',
  `link_updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `link_rel` varchar(255) NOT NULL DEFAULT '',
  `link_notes` mediumtext NOT NULL,
  `link_rss` varchar(255) NOT NULL DEFAULT '',
  `blog_id` bigint(20) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`link_id`),
  KEY `link_visible` (`link_vetoed`)
) DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `postmeta` (
  `postmeta_id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `post_id` bigint(20) unsigned NOT NULL DEFAULT '0',
  `user_id` bigint(20) unsigned NOT NULL DEFAULT '0',
  `meta_key` varchar(255) DEFAULT NULL,
  `meta_value` longtext,
  PRIMARY KEY (`postmeta_id`),
  KEY `post_id` (`post_id`),
  KEY `meta_key` (`meta_key`),
  KEY `user_id` (`user_id`)
) DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `posts` (
  `post_id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) unsigned NOT NULL DEFAULT '0',
  `post_content` longtext NOT NULL,
  `post_title` varchar(200) NOT NULL DEFAULT '',
  `post_excerpt` text NOT NULL,
  `post_status` varchar(20) NOT NULL DEFAULT 'publish',
  `comment_vetoed` varchar(20) NOT NULL DEFAULT '',
  `ping_vetoed` varchar(20) NOT NULL DEFAULT '',
  `post_password` varchar(20) NOT NULL DEFAULT '',
  `post_slug` varchar(200) NOT NULL DEFAULT '',
  `to_ping` text NOT NULL,
  `pinged` text NOT NULL,
  `post_date` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  `post_modified` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `post_parent` bigint(20) unsigned NOT NULL DEFAULT '0',
  `guid` varchar(255) NOT NULL DEFAULT '',
  `menu_order` int(11) unsigned NOT NULL DEFAULT '0',
  `post_mime_type` varchar(100) NOT NULL DEFAULT '',
  `comment_count` bigint(20) unsigned NOT NULL DEFAULT '0',
  `blog_id` bigint(20) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`post_id`),
  KEY `post_slug` (`post_slug`),
  KEY `post_date` (`post_status`,`post_date`,`post_id`),
  KEY `post_parent` (`post_parent`),
  KEY `user_id` (`user_id`)
) DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `users` (
  `user_id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `user_login` char(32) NOT NULL,
  `user_pass` char(32) NOT NULL,
  `user_nicename` varchar(50) NOT NULL,
  `user_date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `user_status` int(11) unsigned NOT NULL DEFAULT '0',
  `site` varchar(20) NOT NULL DEFAULT '',
  PRIMARY KEY (`user_id`),
  UNIQUE KEY `user_login` (`user_login`),
  UNIQUE KEY `site` (`site`)
) DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `usermeta` (
  `usermeta_id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT(20) UNSIGNED NOT NULL DEFAULT '0',
  `meta_tag` VARCHAR(20) NOT NULL,
  `meta_key` VARCHAR(200) NOT NULL,
  `meta_value` LONGTEXT NOT NULL,
  PRIMARY KEY (`usermeta_id`),
  INDEX `user_id` (`user_id`),
  INDEX `meta_key` (`meta_tag`, `meta_key`)
) DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `members` (
  `user_id` bigint(20) unsigned NOT NULL,
  `member_id` bigint(20) unsigned NOT NULL,
  UNIQUE KEY `member_id` (`user_id`,`member_id`)
) DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `sitemeta` (
  `sitemeta_id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) unsigned NOT NULL,
  `meta_key` varchar(255) DEFAULT NULL,
  `meta_value` longtext,
  PRIMARY KEY (`sitemeta_id`),
  KEY `user_id` (`user_id`)
) DEFAULT CHARSET=utf8;
