/*
 Navicat Premium Data Transfer

 Source Server         : 127.0.0.1
 Source Server Type    : MySQL
 Source Server Version : 50733
 Source Host           : 127.0.0.1:3306
 Source Schema         : im

 Target Server Type    : MySQL
 Target Server Version : 50733
 File Encoding         : 65001

 Date: 23/02/2023 03:03:09
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for contact_0
-- ----------------------------
DROP TABLE IF EXISTS `contact_0`;
CREATE TABLE `contact_0` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增id,主键',
  `owner_id` bigint(20) unsigned NOT NULL COMMENT '会话拥有者',
  `peer_id` bigint(20) unsigned NOT NULL COMMENT '联系人（对方用户）',
  `peer_type` tinyint(4) NOT NULL COMMENT '联系人类型，0：普通，100：系统，101：群组',
  `peer_ack` tinyint(4) NOT NULL COMMENT 'peer是否给owner发过消息，0：未发过，1：发过',
  `last_msg_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '双方聊天记录中，最新一次发送的私信id',
  `last_del_msg_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '双方聊天记录中，最后一次删除联系人时的私信id',
  `version_id` bigint(20) unsigned NOT NULL COMMENT '版本号（用于拉取会话框）',
  `sort_key` bigint(20) unsigned NOT NULL COMMENT '会话展示顺序（按顺序展示对话框，也可修改顺序：置顶操作）',
  `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '联系人状态，0：正常，1：被删除',
  `labels` varchar(512) DEFAULT NULL COMMENT '标签,json串',
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `idx_2` (`owner_id`,`peer_id`) USING BTREE,
  KEY `idx_1` (`owner_id`,`version_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='会话表、联系人表（通信双方各有一行记录）';

-- ----------------------------
-- Table structure for message_0
-- ----------------------------
DROP TABLE IF EXISTS `message_0`;
CREATE TABLE `message_0` (
  `msg_id` bigint(20) unsigned NOT NULL COMMENT '私信id',
  `msg_type` tinyint(4) NOT NULL COMMENT '私信类型',
  `session_id` varchar(128) NOT NULL COMMENT '会话id',
  `send_id` bigint(20) unsigned NOT NULL COMMENT '私信发送者',
  `version_id` bigint(20) unsigned NOT NULL COMMENT '版本号（用于拉取消息）',
  `sort_key` bigint(20) unsigned NOT NULL COMMENT '消息展示顺序（按顺序展示消息）',
  `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '私信状态。0：正常，1：被审核删除，2：撤销',
  `content` varchar(2048) NOT NULL COMMENT '私信内容',
  `has_read` tinyint(4) NOT NULL DEFAULT '0' COMMENT '私信接收者是否已读消息。0：未读，1：已读',
  `invisible_list` varchar(2048) DEFAULT NULL COMMENT '消息发出去了，但是对于在列表的用户是不可见的',
  `seq_id` bigint(20) unsigned NOT NULL COMMENT '客户端本地数据库的消息id',
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`msg_id`) USING BTREE,
  KEY `idx_1` (`session_id`,`version_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='一条私信只有一行记录';

SET FOREIGN_KEY_CHECKS = 1;
