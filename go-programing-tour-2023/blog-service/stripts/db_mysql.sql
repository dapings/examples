create database if not exists blog_service
       default character set utf8mb4
       default collate utf8mb4_general_ci;

use blog_service;

drop table if exists `blog_tag`;
create table `blog_tag` (
    `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
    `name` varchar(100) default '' comment '标签名称',
    `created_on` int(10) unsigned default '0' comment '创建时间',
    `modified_on` int(10) unsigned default '0' comment '创建时间',
    `deleted_on` int(10) unsigned default '0' comment '创建时间',
    `is_deleted` tinyint(1) unsigned DEFAULT '0' COMMENT '是否删除 0 未删除、1 已删除',
    `created_by` varchar(128) default '' comment '创建人',
    `modified_by` varchar(128) default '' comment '修改人',
    `status` tinyint(1) unsigned DEFAULT '1' COMMENT '状态 0 禁用、1 启用',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='标签管理';

drop table if exists `blog_article`;
create table `blog_article` (
    `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
    `title` varchar(100) default '' comment '文章标题',
    `desc` varchar(255) default '' comment '文章简述',
    `cover_image_url` varchar(255) default '' comment '封面图片地址',
    `content` longtext comment '文章内容',
    `created_on` int(10) unsigned default '0' comment '创建时间',
    `modified_on` int(10) unsigned default '0' comment '创建时间',
    `deleted_on` int(10) unsigned default '0' comment '创建时间',
    `is_deleted` tinyint(1) unsigned DEFAULT '0' COMMENT '是否删除 0 未删除、1 已删除',
    `created_by` varchar(128) default '' comment '创建人',
    `modified_by` varchar(128) default '' comment '修改人',
    `status` tinyint(1) unsigned DEFAULT '1' COMMENT '状态 0 禁用、1 启用',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文章管理';

drop table if exists `blog_article_tag`;
create table `blog_article_tag` (
    `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
    `article_id` int(10) not null comment '文章ID',
    `tag_id` int(10) unsigned default '0' comment '标签ID',
    `created_on` int(10) unsigned default '0' comment '创建时间',
    `modified_on` int(10) unsigned default '0' comment '创建时间',
    `deleted_on` int(10) unsigned default '0' comment '创建时间',
    `is_deleted` tinyint(1) unsigned DEFAULT '0' COMMENT '是否删除 0 未删除、1 已删除',
    `created_by` varchar(128) default '' comment '创建人',
    `modified_by` varchar(128) default '' comment '修改人',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文章标签关联';