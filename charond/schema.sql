CREATE SCHEMA IF NOT EXISTS charon;

CREATE TABLE IF NOT EXISTS charon.user (
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

  CONSTRAINT "charon.user_pkey" PRIMARY KEY (id),
  CONSTRAINT "charon.user_username_key" UNIQUE (username),
  CONSTRAINT "charon.user_created_by_fkey" FOREIGN KEY (created_by) REFERENCES charon.user (id),
  CONSTRAINT "charon.user_updated_by_fkey" FOREIGN KEY (updated_by) REFERENCES charon.user (id)
);

CREATE TABLE IF NOT EXISTS charon.group (
  id         SERIAL,
  name       TEXT                      NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
  created_by INTEGER,
  updated_at TIMESTAMPTZ,
  updated_by INTEGER,

  CONSTRAINT "charon.group_pkey" PRIMARY KEY (id),
  CONSTRAINT "charon.group_created_by_fkey" FOREIGN KEY (created_by) REFERENCES charon.user (id),
  CONSTRAINT "charon.group_updated_by_fkey" FOREIGN KEY (updated_by) REFERENCES charon.user (id)
);

CREATE TABLE IF NOT EXISTS charon.user_groups (
  user_id    INTEGER                   NOT NULL,
  group_id   INTEGER                   NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
  created_by INTEGER,

  CONSTRAINT "charon.user_groups_user_id_fkey" FOREIGN KEY (user_id) REFERENCES charon.user (id),
  CONSTRAINT "charon.user_groups_group_id_fkey" FOREIGN KEY (group_id) REFERENCES charon.group (id),
  CONSTRAINT "charon.user_groups_created_by_fkey" FOREIGN KEY (created_by) REFERENCES charon.user (id)
);

CREATE TABLE IF NOT EXISTS charon.permission (
  id           SERIAL,
  subsystem_id INTEGER,
  subsystem    TEXT                      NOT NULL,
  module       TEXT                      NOT NULL,
  action       TEXT                      NOT NULL,
  created_at   TIMESTAMPTZ DEFAULT NOW() NOT NULL,
  created_by   INTEGER,

  CONSTRAINT "charon.permission_pkey" PRIMARY KEY (id),
  CONSTRAINT "charon.permission_subsystem_module_action_key" UNIQUE (subsystem, module, action),
  CONSTRAINT "charon.permission_subsystem_id_fkey" FOREIGN KEY (subsystem_id) REFERENCES charon.subsystem (id),
  CONSTRAINT "charon.permission_created_by_fkey" FOREIGN KEY (created_by) REFERENCES charon.user (id)
);

CREATE TABLE IF NOT EXISTS charon.user_permissions (
  user_id       INTEGER                   NOT NULL,
  permission_id INTEGER                   NOT NULL,
  created_at    TIMESTAMPTZ DEFAULT NOW() NOT NULL,
  created_by    INTEGER,

  CONSTRAINT "charon.user_permissions_user_id_fkey" FOREIGN KEY (user_id) REFERENCES charon.user (id),
  CONSTRAINT "charon.user_permissions_permission_id_fkey" FOREIGN KEY (permission_id) REFERENCES charon.permission (id),
  CONSTRAINT "charon.user_permissions_created_by_fkey" FOREIGN KEY (created_by) REFERENCES charon.user (id)
);


CREATE TABLE IF NOT EXISTS charon.group_permissions (
  group_id      INTEGER                   NOT NULL,
  permission_id INTEGER                   NOT NULL,
  created_at    TIMESTAMPTZ DEFAULT NOW() NOT NULL,
  created_by    INTEGER,

  CONSTRAINT "charon.group_permissions_group_id_fkey" FOREIGN KEY (group_id) REFERENCES charon.group (id),
  CONSTRAINT "charon.group_permissions_permission_id_fkey" FOREIGN KEY (permission_id) REFERENCES charon.permission (id),
  CONSTRAINT "charon.group_permissions_created_by_fkey" FOREIGN KEY (created_by) REFERENCES charon.user (id)
);