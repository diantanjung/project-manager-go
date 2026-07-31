CREATE TYPE user_role AS ENUM ('admin', 'productOwner', 'projectManager', 'teamMember');
CREATE TYPE team_member_role AS ENUM ('owner', 'admin', 'member');
CREATE TYPE project_status AS ENUM ('planning', 'active', 'paused', 'completed', 'archived');
CREATE TYPE task_status AS ENUM ('backlog', 'todo', 'in_progress', 'review', 'done');
CREATE TYPE task_priority AS ENUM ('low', 'medium', 'high', 'urgent');
CREATE TYPE notification_type AS ENUM ('task_assigned', 'mention', 'task_due', 'project_update', 'system_alert');
CREATE TYPE webhook_event AS ENUM ('task.created', 'task.updated', 'task.completed', 'comment.created', 'project.updated');
CREATE TYPE webhook_status AS ENUM ('pending', 'delivered', 'failed');
CREATE TYPE export_status AS ENUM ('pending', 'processing', 'completed', 'failed');
