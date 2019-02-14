CREATE TABLE IF NOT EXISTS user (
  user_id      varchar(50)  NOT NULL,
  username     varchar(50)  NOT NULL,
  email        varchar(50)  NOT NULL,
  phone_number varchar(50)           DEFAULT NULL,
  description  varchar(255) NOT NULL,
  password     varchar(255) NOT NULL,
  status       varchar(50)  NOT NULL,
  create_time  timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_time  timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP,
  status_time  timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP,
  extra        json                  DEFAULT NULL,
  PRIMARY KEY (user_id)
);

CREATE INDEX user_email_idx
  ON user (email);
CREATE INDEX user_phone_number_idx
  ON user (phone_number);
CREATE INDEX user_status_idx
  ON user (status);
CREATE INDEX user_username_idx
  ON user (username);
CREATE INDEX user_create_time_idx
  ON user (create_time);

CREATE TABLE IF NOT EXISTS `group` (
  group_id         varchar(50)  NOT NULL,
  parent_group_id  varchar(50)  NOT NULL,
  group_path       varchar(255) NOT NULL,
  group_name       varchar(50)  NOT NULL,
  description      varchar(255)          DEFAULT NULL,
  status           varchar(50)  NOT NULL,
  create_time      timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_time      timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP,
  status_time      timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP,
  extra            json                  DEFAULT NULL,
  group_path_level int(11)      NOT NULL,
  PRIMARY KEY (group_id)
);
CREATE INDEX group_parent_group_id_idx
  ON `group` (parent_group_id);
CREATE INDEX group_group_path_idx
  ON `group` (group_path);
CREATE INDEX group_group_name_idx
  ON `group` (group_name);
CREATE INDEX group_status_idx
  ON `group` (status);
CREATE INDEX group_group_path_level_idx
  ON `group` (group_path_level);
CREATE INDEX group_create_time_idx
  ON `group` (create_time);

CREATE TABLE IF NOT EXISTS user_group_binding (
  id          varchar(50) NOT NULL,
  user_id     varchar(50) NOT NULL,
  group_id    varchar(50) NOT NULL,
  create_time timestamp   NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id)
);
CREATE INDEX user_group_binding_user_id_idx
  ON user_group_binding (user_id);
CREATE INDEX user_group_binding_group_id_idx
  ON user_group_binding (group_id);