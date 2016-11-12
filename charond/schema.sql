CREATE SCHEMA IF NOT EXISTS charon;

CREATE TABLE IF NOT EXISTS charonrpc.user (
  id                 SERIAL,
  password           TEXT                      NOT NULL,
  username           TEXT                      NOT NULL,
  first_name         TEXT                      NOT NULL,
  last_name          TEXT                      NOT NULL,
  is_superuser       BOOLEAN                   NOT NULL,
  is_active          BOOLEAN                   NOT NULL,
  is_staff           BOOLEAN                   NOT NULL,
  is_confirmed       BOOLEAN                   NOT NULL,
  confirmation_token TEXT                      NOT NULL,
  last_login_at      TIMESTAMPTZ,
  created_at         TIMESTAMPTZ DEFAULT NOW() NOT NULL,
  created_by         INTEGER,
  updated_at         TIMESTAMPTZ,
  updated_by         INTEGER,

  CONSTRAINT "charonrpc.user_pkey" PRIMARY KEY (id),
  CONSTRAINT "charonrpc.user_username_key" UNIQUE (username),
  CONSTRAINT "charonrpc.user_created_by_fkey" FOREIGN KEY (created_by) REFERENCES charonrpc.user (id),
  CONSTRAINT "charonrpc.user_updated_by_fkey" FOREIGN KEY (updated_by) REFERENCES charonrpc.user (id)
);

CREATE TABLE IF NOT EXISTS charonrpc.group (
  id         SERIAL,
  name       TEXT                      NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
  created_by INTEGER,
  updated_at TIMESTAMPTZ,
  updated_by INTEGER,

  CONSTRAINT "charonrpc.group_pkey" PRIMARY KEY (id),
  CONSTRAINT "charonrpc.group_created_by_fkey" FOREIGN KEY (created_by) REFERENCES charonrpc.user (id),
  CONSTRAINT "charonrpc.group_updated_by_fkey" FOREIGN KEY (updated_by) REFERENCES charonrpc.user (id)
);

CREATE TABLE IF NOT EXISTS charonrpc.user_groups (
  user_id    INTEGER                   NOT NULL,
  group_id   INTEGER                   NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
  created_by INTEGER,

  CONSTRAINT "charonrpc.user_groups_user_id_fkey" FOREIGN KEY (user_id) REFERENCES charonrpc.user (id),
  CONSTRAINT "charonrpc.user_groups_group_id_fkey" FOREIGN KEY (group_id) REFERENCES charonrpc.group (id),
  CONSTRAINT "charonrpc.user_groups_created_by_fkey" FOREIGN KEY (created_by) REFERENCES charonrpc.user (id)
);

CREATE TABLE IF NOT EXISTS charonrpc.permission (
  id           SERIAL,
  subsystem_id INTEGER,
  subsystem    TEXT                      NOT NULL,
  module       TEXT                      NOT NULL,
  action       TEXT                      NOT NULL,
  created_at   TIMESTAMPTZ DEFAULT NOW() NOT NULL,
  created_by   INTEGER,

  CONSTRAINT "charonrpc.permission_pkey" PRIMARY KEY (id),
  CONSTRAINT "charonrpc.permission_subsystem_module_action_key" UNIQUE (subsystem, module, action),
  CONSTRAINT "charonrpc.permission_subsystem_id_fkey" FOREIGN KEY (subsystem_id) REFERENCES charonrpc.subsystem (id),
  CONSTRAINT "charonrpc.permission_created_by_fkey" FOREIGN KEY (created_by) REFERENCES charonrpc.user (id)
);

CREATE TABLE IF NOT EXISTS charonrpc.user_permissions (
  user_id       INTEGER                   NOT NULL,
  permission_id INTEGER                   NOT NULL,
  created_at    TIMESTAMPTZ DEFAULT NOW() NOT NULL,
  created_by    INTEGER,

  CONSTRAINT "charonrpc.user_permissions_user_id_fkey" FOREIGN KEY (user_id) REFERENCES charonrpc.user (id),
  CONSTRAINT "charonrpc.user_permissions_permission_id_fkey" FOREIGN KEY (permission_id) REFERENCES charonrpc.permission (id),
  CONSTRAINT "charonrpc.user_permissions_created_by_fkey" FOREIGN KEY (created_by) REFERENCES charonrpc.user (id)
);


CREATE TABLE IF NOT EXISTS charonrpc.group_permissions (
  group_id      INTEGER                   NOT NULL,
  permission_id INTEGER                   NOT NULL,
  created_at    TIMESTAMPTZ DEFAULT NOW() NOT NULL,
  created_by    INTEGER,

  CONSTRAINT "charonrpc.group_permissions_group_id_fkey" FOREIGN KEY (group_id) REFERENCES charonrpc.group (id),
  CONSTRAINT "charonrpc.group_permissions_permission_id_fkey" FOREIGN KEY (permission_id) REFERENCES charonrpc.permission (id),
  CONSTRAINT "charonrpc.group_permissions_created_by_fkey" FOREIGN KEY (created_by) REFERENCES charonrpc.user (id)
);