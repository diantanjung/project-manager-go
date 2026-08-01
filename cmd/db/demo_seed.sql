BEGIN;
-- pass: password123
INSERT INTO users (id, name, email, password, created_at, updated_at, avatar_storage_key, role)
VALUES
	(27, 'Admin User', 'admin@example.com', '$2y$12$XhudwBrQ/5V3PCLNz2TZ4ecmgPwcmwUjYE5CcgG1bdbEdZ3861miO', '2026-02-20 04:56:53.745518+00', '2026-02-20 04:56:53.745518+00', 'https://api.dicebear.com/7.x/avataaars/svg?seed=admin', 'admin'),
	(28, 'John Product Owner', 'john@example.com', '$2y$12$XhudwBrQ/5V3PCLNz2TZ4ecmgPwcmwUjYE5CcgG1bdbEdZ3861miO', '2026-02-20 04:56:53.745518+00', '2026-02-20 04:56:53.745518+00', 'https://api.dicebear.com/7.x/avataaars/svg?seed=john', 'productOwner'),
	(29, 'Jane Manager', 'jane@example.com', '$2y$12$XhudwBrQ/5V3PCLNz2TZ4ecmgPwcmwUjYE5CcgG1bdbEdZ3861miO', '2026-02-20 04:56:53.745518+00', '2026-02-20 04:56:53.745518+00', 'https://api.dicebear.com/7.x/avataaars/svg?seed=jane', 'projectManager'),
	(30, 'Bob Manager', 'bob@example.com', '$2y$12$XhudwBrQ/5V3PCLNz2TZ4ecmgPwcmwUjYE5CcgG1bdbEdZ3861miO', '2026-02-20 04:56:53.745518+00', '2026-02-20 04:56:53.745518+00', 'https://api.dicebear.com/7.x/avataaars/svg?seed=bob', 'projectManager'),
	(31, 'Alice Developer', 'alice@example.com', '$2y$12$XhudwBrQ/5V3PCLNz2TZ4ecmgPwcmwUjYE5CcgG1bdbEdZ3861miO', '2026-02-20 04:56:53.745518+00', '2026-02-20 04:56:53.745518+00', 'https://api.dicebear.com/7.x/avataaars/svg?seed=alice', 'teamMember'),
	(32, 'Charlie Designer', 'charlie@example.com', '$2y$12$XhudwBrQ/5V3PCLNz2TZ4ecmgPwcmwUjYE5CcgG1bdbEdZ3861miO', '2026-02-20 04:56:53.745518+00', '2026-02-20 04:56:53.745518+00', 'https://api.dicebear.com/7.x/avataaars/svg?seed=charlie', 'teamMember'),
	(33, 'Dave Developer', 'dave@example.com', '$2y$12$XhudwBrQ/5V3PCLNz2TZ4ecmgPwcmwUjYE5CcgG1bdbEdZ3861miO', '2026-02-20 04:56:53.745518+00', '2026-02-20 04:56:53.745518+00', 'https://api.dicebear.com/7.x/avataaars/svg?seed=dave', 'teamMember'),
	(34, 'Test User', 'unique@example.com', '$2y$12$XhudwBrQ/5V3PCLNz2TZ4ecmgPwcmwUjYE5CcgG1bdbEdZ3861miO', '2026-03-05 17:18:19.064589+00', '2026-03-05 17:18:19.064589+00', NULL, 'teamMember')
ON CONFLICT (email) DO UPDATE SET
	name = EXCLUDED.name,
	password = EXCLUDED.password,
	avatar_storage_key = EXCLUDED.avatar_storage_key,
	role = EXCLUDED.role,
	updated_at = EXCLUDED.updated_at;

INSERT INTO teams (id, name, description, created_at, updated_at)
VALUES
	(7, 'Engineering', 'Core engineering team responsible for product development', '2026-02-20 04:56:53.749393+00', '2026-02-20 04:56:53.749393+00'),
	(8, 'Design', 'UX/UI design team creating beautiful user experiences', '2026-02-20 04:56:53.749393+00', '2026-02-20 04:56:53.749393+00'),
	(9, 'QA', 'Quality assurance team ensuring product quality', '2026-02-20 04:56:53.749393+00', '2026-02-20 04:56:53.749393+00')
ON CONFLICT (id) DO UPDATE SET
	name = EXCLUDED.name,
	description = EXCLUDED.description,
	updated_at = EXCLUDED.updated_at;

INSERT INTO team_members (id, team_id, user_id, role, joined_at)
VALUES
	(17, 7, 29, 'owner', '2026-02-20 04:56:53.750688+00'),
	(18, 7, 31, 'member', '2026-02-20 04:56:53.750688+00'),
	(19, 7, 33, 'member', '2026-02-20 04:56:53.750688+00'),
	(20, 8, 30, 'owner', '2026-02-20 04:56:53.750688+00'),
	(21, 8, 32, 'member', '2026-02-20 04:56:53.750688+00'),
	(22, 9, 29, 'admin', '2026-02-20 04:56:53.750688+00'),
	(23, 9, 33, 'member', '2026-02-20 04:56:53.750688+00')
ON CONFLICT (team_id, user_id) DO UPDATE SET
	role = EXCLUDED.role,
	joined_at = EXCLUDED.joined_at;

INSERT INTO projects (id, name, description, team_id, owner_id, created_at, updated_at)
VALUES
	(7, 'Website Redesign', 'Complete overhaul of the company website with modern design and improved UX', 7, 28, '2026-02-20 04:56:53.752588+00', '2026-02-20 04:56:53.752588+00'),
	(8, 'Mobile App v2.0', 'Major update to the mobile application with new features and performance improvements', 7, 28, '2026-02-20 04:56:53.752588+00', '2026-02-20 04:56:53.752588+00'),
	(9, 'Brand Guidelines', 'Comprehensive brand identity guidelines and design system documentation', 8, 28, '2026-02-20 04:56:53.752588+00', '2026-02-20 04:56:53.752588+00')
ON CONFLICT (id) DO UPDATE SET
	name = EXCLUDED.name,
	description = EXCLUDED.description,
	team_id = EXCLUDED.team_id,
	owner_id = EXCLUDED.owner_id,
	updated_at = EXCLUDED.updated_at;

INSERT INTO project_teams (id, project_id, team_id, assigned_at)
VALUES
	(6, 7, 7, '2026-02-20 04:56:53.754131+00'),
	(7, 7, 8, '2026-02-20 04:56:53.754131+00'),
	(8, 8, 7, '2026-02-20 04:56:53.754131+00'),
	(9, 8, 9, '2026-02-20 04:56:53.754131+00'),
	(10, 9, 8, '2026-02-20 04:56:53.754131+00')
ON CONFLICT (project_id, team_id) DO UPDATE SET
	assigned_at = EXCLUDED.assigned_at;

INSERT INTO tasks (id, title, description, status, priority, project_id, creator_id, assignee_id, due_date, position, created_at, updated_at)
VALUES
	(30, 'Design homepage mockup', 'Create high-fidelity mockup for the new homepage including hero section, features, and footer', 'done', 'high', 7, 29, 32, '2026-02-05', 0, '2026-02-20 04:56:53.755281+00', '2026-02-20 04:56:53.755281+00'),
	(31, 'Implement responsive navigation', 'Build responsive navbar with mobile hamburger menu and smooth transitions', 'in_progress', 'high', 7, 29, 31, '2026-02-08', 1, '2026-02-20 04:56:53.755281+00', '2026-02-20 04:56:53.755281+00'),
	(32, 'Setup CI/CD pipeline', 'Configure GitHub Actions for automated testing and deployment to staging/production', 'review', 'medium', 7, 29, 33, '2026-02-10', 0, '2026-02-20 04:56:53.755281+00', '2026-02-20 04:56:53.755281+00'),
	(33, 'Optimize images and assets', 'Compress images, implement lazy loading, and setup CDN for static assets', 'todo', 'medium', 7, 29, 31, '2026-02-12', 0, '2026-02-20 04:56:53.755281+00', '2026-02-20 04:56:53.755281+00'),
	(34, 'Write unit tests for components', 'Add comprehensive unit tests for all React components using Jest and Testing Library', 'backlog', 'low', 7, 29, 33, '2026-02-15', 0, '2026-02-20 04:56:53.755281+00', '2026-02-20 04:56:53.755281+00'),
	(35, 'Implement push notifications', 'Setup Firebase Cloud Messaging for iOS and Android push notifications', 'in_progress', 'urgent', 8, 29, 31, '2026-02-06', 0, '2026-02-20 04:56:53.755281+00', '2026-02-20 04:56:53.755281+00'),
	(36, 'Refactor authentication flow', 'Improve the login/signup flow with biometric authentication support', 'todo', 'high', 8, 29, 33, '2026-02-09', 0, '2026-02-20 04:56:53.755281+00', '2026-02-20 04:56:53.755281+00'),
	(37, 'Add offline mode support', 'Implement local storage caching and sync mechanism for offline usage', 'backlog', 'medium', 8, 29, 31, '2026-02-20', 1, '2026-02-20 04:56:53.755281+00', '2026-02-20 04:56:53.755281+00'),
	(38, 'Define color palette', 'Establish primary, secondary, and accent colors with accessibility considerations', 'done', 'high', 9, 30, 32, '2026-01-28', 0, '2026-02-20 04:56:53.755281+00', '2026-02-20 04:56:53.755281+00'),
	(39, 'Create typography guide', 'Select and document font families, sizes, weights, and usage guidelines', 'in_progress', 'high', 9, 30, 32, '2026-02-04', 0, '2026-02-20 04:56:53.755281+00', '2026-02-20 04:56:53.755281+00'),
	(40, 'Design icon library', 'Create consistent iconography set for use across all platforms', 'todo', 'medium', 9, 30, 32, '2026-02-14', 0, '2026-02-20 04:56:53.755281+00', '2026-02-20 04:56:53.755281+00')
ON CONFLICT (id) DO UPDATE SET
	title = EXCLUDED.title,
	description = EXCLUDED.description,
	status = EXCLUDED.status,
	priority = EXCLUDED.priority,
	project_id = EXCLUDED.project_id,
	creator_id = EXCLUDED.creator_id,
	assignee_id = EXCLUDED.assignee_id,
	due_date = EXCLUDED.due_date,
	position = EXCLUDED.position,
	updated_at = EXCLUDED.updated_at;

INSERT INTO task_assignments (task_id, user_id, assigned_at)
SELECT id, assignee_id, created_at FROM tasks WHERE assignee_id IS NOT NULL
ON CONFLICT (task_id, user_id) DO UPDATE SET
	assigned_at = EXCLUDED.assigned_at;

INSERT INTO comments (id, content, task_id, author_id, created_at, updated_at)
VALUES
	(22, 'I''ve completed the initial wireframes. Please review when you get a chance!', 30, 32, '2026-02-20 04:56:53.760327+00', '2026-02-20 04:56:53.760327+00'),
	(23, 'Looks great! I love the hero section. Can we add more whitespace around the CTA button?', 30, 29, '2026-02-20 04:56:53.760327+00', '2026-02-20 04:56:53.760327+00'),
	(24, 'Done! Updated the mockup with extra padding. Ready for final review.', 30, 32, '2026-02-20 04:56:53.760327+00', '2026-02-20 04:56:53.760327+00'),
	(25, 'Started working on this. Should we use a slide-in or dropdown menu for mobile?', 31, 31, '2026-02-20 04:56:53.760327+00', '2026-02-20 04:56:53.760327+00'),
	(26, 'Let''s go with slide-in from the left. It feels more modern and works better with our layout.', 31, 29, '2026-02-20 04:56:53.760327+00', '2026-02-20 04:56:53.760327+00'),
	(27, 'PR is ready for review. I''ve added deployment to staging on merge to develop branch.', 32, 33, '2026-02-20 04:56:53.760327+00', '2026-02-20 04:56:53.760327+00'),
	(28, 'Need access to Firebase console to proceed. @admin can you add me?', 35, 31, '2026-02-20 04:56:53.760327+00', '2026-02-20 04:56:53.760327+00'),
	(29, 'Done! You should have admin access now.', 35, 27, '2026-02-20 04:56:53.760327+00', '2026-02-20 04:56:53.760327+00'),
	(30, 'I''m thinking Inter for body text and Playfair Display for headings. Thoughts?', 39, 32, '2026-02-20 04:56:53.760327+00', '2026-02-20 04:56:53.760327+00'),
	(31, 'Love the combination! Inter is super readable. Let''s make sure we have the full weight range.', 39, 30, '2026-02-20 04:56:53.760327+00', '2026-02-20 04:56:53.760327+00')
ON CONFLICT (id) DO UPDATE SET
	content = EXCLUDED.content,
	task_id = EXCLUDED.task_id,
	author_id = EXCLUDED.author_id,
	updated_at = EXCLUDED.updated_at;

INSERT INTO attachments (id, task_id, uploader_id, file_name, original_name, storage_key, mime_type, file_size, created_at, updated_at)
VALUES
	(11, 30, 32, 'homepage-mockup-v1.fig', 'homepage-mockup-v1.fig', 'attachments/homepage-mockup-v1.fig', 'application/figma', 2457600, '2026-02-20 04:56:53.761304+00', '2026-02-20 04:56:53.761304+00'),
	(12, 30, 32, 'homepage-mockup-v2.fig', 'homepage-mockup-v2.fig', 'attachments/homepage-mockup-v2.fig', 'application/figma', 2867200, '2026-02-20 04:56:53.761304+00', '2026-02-20 04:56:53.761304+00'),
	(13, 32, 33, 'ci-cd-diagram.png', 'ci-cd-diagram.png', 'attachments/ci-cd-diagram.png', 'image/png', 156800, '2026-02-20 04:56:53.761304+00', '2026-02-20 04:56:53.761304+00'),
	(14, 38, 32, 'color-palette.pdf', 'color-palette.pdf', 'attachments/color-palette.pdf', 'application/pdf', 524288, '2026-02-20 04:56:53.761304+00', '2026-02-20 04:56:53.761304+00'),
	(15, 39, 32, 'typography-samples.pdf', 'typography-samples.pdf', 'attachments/typography-samples.pdf', 'application/pdf', 1048576, '2026-02-20 04:56:53.761304+00', '2026-02-20 04:56:53.761304+00')
ON CONFLICT (id) DO UPDATE SET
	task_id = EXCLUDED.task_id,
	uploader_id = EXCLUDED.uploader_id,
	file_name = EXCLUDED.file_name,
	original_name = EXCLUDED.original_name,
	storage_key = EXCLUDED.storage_key,
	mime_type = EXCLUDED.mime_type,
	file_size = EXCLUDED.file_size,
	updated_at = EXCLUDED.updated_at;

INSERT INTO refresh_tokens (id, user_id, token, is_revoked, expires_at, created_at)
VALUES
	(130, 27, '2aa903bb372d9faee265b68b2a1e0e6456950e3e283a3513f7375ae2d029df9b', true, '2026-04-03 02:39:11.036+00', '2026-03-03 18:39:11.036778+00'),
	(131, 27, '963dd79fa4c25250aa518b2a1ce1fcd370d36390c127d039de1ea75da4895444', true, '2026-04-03 06:40:42.095+00', '2026-03-03 22:40:42.095646+00'),
	(132, 27, '320ced6053ff68b7af186a4a5538bff9426ac564a16a3ad57593352313a60ead', true, '2026-04-03 06:40:43.516+00', '2026-03-03 22:40:43.517168+00'),
	(133, 27, '46ec4fd71c2dfccc82be655a23c57240e5006307e4aa468ae8b48b0a6e949239', true, '2026-04-03 06:40:44.115+00', '2026-03-03 22:40:44.116039+00'),
	(134, 27, 'c0313db8a98cdf3c3e980c982e18bff8bf388e9e427f4c00133485d1220e144a', true, '2026-04-03 06:40:48.504+00', '2026-03-03 22:40:48.504627+00'),
	(135, 27, '861cc879e14fd8ca32c361baa84b630b2ccdee836dd79f7c3fd018328161eafa', true, '2026-04-03 08:38:38.601+00', '2026-03-04 00:38:38.602033+00'),
	(136, 27, 'ae597cd76e1e184f8610504a29ee4cda027a54481ed89c97cd38805d8e016de7', true, '2026-04-04 09:22:38.315+00', '2026-03-05 01:22:38.31586+00'),
	(137, 27, '69b39a56d9a73bc9e23f3efa239e52d739a68c35b9eca3671c8c3e175e459a7d', true, '2026-04-04 13:29:28.405+00', '2026-03-05 05:29:28.405712+00'),
	(138, 27, 'ec0fe2c50912cf2572015c1cadaf1f7d39648fe960a8fafe1a045b409067e94a', true, '2026-04-05 01:01:15.881+00', '2026-03-05 17:01:15.881788+00'),
	(139, 34, '3e8d29672a9414409829e8a8be1e42b4bb3cc82cd46a1dcfbba0fa3a22503685', true, '2026-04-05 01:18:59.869+00', '2026-03-05 17:18:59.869756+00'),
	(140, 27, 'd2e58acec39cb0be4853ad0cbc3997a33908e0307403f339efa4a71d399e3ed9', true, '2026-04-05 01:46:41.297+00', '2026-03-05 17:46:41.297242+00'),
	(141, 27, 'bcf7fdb5bab29ac252b44a479ae8c761bd996b8f36c026408aeaa50525998e47', true, '2026-04-05 01:49:04.332+00', '2026-03-05 17:49:04.33274+00'),
	(142, 27, 'c0389ed84de2c698d133e310d60f3fb810e52738a32101d6a0a6ddc0f59b11dc', true, '2026-04-05 01:54:28.046+00', '2026-03-05 17:54:28.047128+00'),
	(143, 34, 'a5c023bbb0bcc1fdcd97db06f9414560117433b3d5925837f8a4cedf6f98ad4d', false, '2026-04-05 03:07:02.918+00', '2026-03-05 19:07:02.918252+00'),
	(144, 27, '03fb7eb5b4e0fe4b397bd90f0e0fd6fbd6f1b20bf2ac5e7f06bc80c19650f948', true, '2026-04-05 03:07:02.979+00', '2026-03-05 19:07:02.979257+00'),
	(145, 27, '55cac378ac7519f5cf94a0f0cb16b11579c2792051165fdc9c5c8c1b92ec83a1', true, '2026-04-05 08:39:46.515+00', '2026-03-06 00:39:46.515696+00'),
	(146, 27, '22782dd4506a026f3ff5fc202a667ad36bf0db973dd2438bc593f1ba08e60bb2', true, '2026-04-06 03:39:08.862+00', '2026-03-06 19:39:08.862586+00'),
	(147, 27, '33c86e3b48d1f4ca7516d99cdacfefaf93c7f1ed852f325ad1076e5ef54d30cc', false, '2026-04-06 03:42:42.27+00', '2026-03-06 19:42:42.271538+00'),
	(148, 27, 'aea506edcb978fd10382eec44029e419147128b45e26d203e2f2e94ad2802e22', true, '2026-08-07 07:21:34.98+00', '2026-07-08 00:21:34.981082+00'),
	(149, 27, '89ce7d90e23f53cb493d23c53078b46a5cf5a71ca67e57d87aafb49fc08df20a', true, '2026-08-07 08:16:45.352+00', '2026-07-08 01:16:45.352948+00'),
	(150, 27, '1eda3fc24a3e63f129a79adeaebbb13eb9f824162089523ecd88a14deea8d1ff', true, '2026-08-07 10:44:04.595+00', '2026-07-08 03:44:04.595862+00'),
	(151, 27, '7c22dd7d5c9d30707388c36e4df9803df75bccf24f4c88118a0d10d1665c4755', true, '2026-08-07 11:49:34.364+00', '2026-07-08 04:49:34.364415+00'),
	(152, 27, '3df4b3d2f0a3f4d5cef8e5e51d2f5e7846d73b371e862aa582b039c95cd2c2a6', true, '2026-08-07 13:21:34.869+00', '2026-07-08 06:21:34.869993+00'),
	(153, 27, '68de176ed586bcbd202c5d5401ce7e64afe8d3bf452c2cd78c15cc11a85a6708', true, '2026-08-07 22:33:12.057+00', '2026-07-08 15:33:12.058052+00'),
	(154, 27, '04046f244dd1e176e6acd33094e7e18f06a6852935a0ffd13a712afb4c4d6c1a', true, '2026-08-08 00:01:21.952+00', '2026-07-08 17:01:21.952427+00'),
	(155, 27, 'e16aec63ee86e551cd776ce5e5247997a9336333e1e95fe48497c32bcbd7df51', true, '2026-08-08 01:50:45.607+00', '2026-07-08 18:50:45.608052+00'),
	(156, 27, '6d1fab3a9abf7d6910b07d554a1d2e406f2b1eae911a2c5464280906798030d3', true, '2026-08-08 02:44:07.816+00', '2026-07-08 19:44:07.817753+00'),
	(157, 27, '368f767d907a41cdf2cd6125e22bfaebbcb86a8ce150d28c4af4a0658c27c6dc', true, '2026-08-08 03:01:47.665+00', '2026-07-08 20:01:47.66588+00'),
	(158, 27, '2877b1a21cde9545553b0353ad0031bc58982dcc2889ddd03dfc0084e87cb7c5', true, '2026-08-08 03:12:39.554+00', '2026-07-08 20:12:39.554152+00'),
	(159, 27, 'a52390686353122edb0d4d2e7fe92e42e8d9216f78cdb1aefb0096a5da77eb58', true, '2026-08-08 04:18:35.782+00', '2026-07-08 21:18:35.782476+00'),
	(160, 27, '2aa1116bdb499daeed389016ad8bd9d2918c885de2ba7e908c3cc3f1515becec', true, '2026-08-08 05:37:35.53+00', '2026-07-08 22:37:35.530502+00'),
	(161, 27, '5eb4941a1101103062ea4883bb7b3e2d83520266dc3766b8d98642ff52c8bf67', true, '2026-08-08 07:21:02.997+00', '2026-07-09 00:21:02.997425+00'),
	(162, 27, '31be082f5f862e309b088ccfdb9094b379a77d84762cc992d2d7ba02310f669f', true, '2026-08-08 07:32:40.083+00', '2026-07-09 00:32:40.083643+00'),
	(163, 27, 'fa324b488d6115b4ba77ea813f117bb57146c5dbdf57a5431a4436cc2c2082c8', true, '2026-08-08 08:23:50.229+00', '2026-07-09 01:23:50.229665+00'),
	(164, 27, 'ab4db1c9875d95ec880fc94a41b7eb916e1178eeacaaf84ee3a5cfe77eede88f', true, '2026-08-08 10:35:23.076+00', '2026-07-09 03:35:23.07673+00'),
	(165, 27, '48964a28e4572d7f12120e3fd74c7870a175d07d4befdc42c94d11f011b71b70', true, '2026-08-08 21:41:14.246+00', '2026-07-09 14:41:14.246548+00'),
	(166, 27, '6d8f5eaa8c023bb94dcfd1fd0c4f31bb8fa032666b03525545501ea8562c0385', true, '2026-08-09 02:41:32.585+00', '2026-07-09 19:41:32.585395+00'),
	(167, 27, '8e7e83bf83ff03835b3d22247c4ac04d00ee09aa1738b5d3f424eb58262a14f1', true, '2026-08-09 04:14:38.065+00', '2026-07-09 21:14:38.066071+00'),
	(168, 27, 'a13fd55223d7de1b6c27c946baee8d268c17ca92feb62646cd80badd9e50203c', true, '2026-08-09 06:23:55.406+00', '2026-07-09 23:23:55.407023+00'),
	(169, 27, '789b43ea118d385078552f9c839a66be8b03c793b77aa40d691997426e9d4964', true, '2026-08-09 08:08:32.611+00', '2026-07-10 01:08:32.6119+00'),
	(170, 27, 'b9a801bed8224a46a824bda321f93ce1a66e6f78589c1fe3635c8cbee646f9b0', true, '2026-08-09 10:46:49.372+00', '2026-07-10 03:46:49.372565+00'),
	(171, 27, 'c7aaeefef65253fd84f33e2ed529e97a61c5e618ed8e0c564d03d1d13a46b011', true, '2026-08-09 12:23:44.243+00', '2026-07-10 05:23:44.243835+00'),
	(172, 27, '07b1d447222676c70d0037aa22b9a6579ec7468a6386d7a1e1c97d5e94c12433', true, '2026-08-09 13:52:30.512+00', '2026-07-10 06:52:30.51352+00'),
	(173, 27, 'c90f1e04e42845bb4ef4465fa3412feb549e0278bbae0e1539f32d6748bddf59', true, '2026-08-09 15:17:23.753+00', '2026-07-10 08:17:23.753449+00'),
	(174, 27, '260bd61ea4442f7c25e0ec6dc4e8addf53be2495653d34d0c8053064feba9052', true, '2026-08-09 15:25:13.573+00', '2026-07-10 08:25:13.574355+00'),
	(175, 27, 'ddd58f7f3262393fc5a8bc8e6924a53cd4dab1018f2ce8c8852bab666bb0c9a8', true, '2026-08-09 15:50:08.102+00', '2026-07-10 08:50:08.10288+00'),
	(177, 27, '57a2f5b8df0a1ea772089db4fa2908cf9ae9cb4923fb012040ea79fa0e0bb3c5', false, '2026-08-11 02:52:29.937+00', '2026-07-11 19:52:29.938436+00'),
	(178, 27, '2c636528c0b3ef4f9dedc47403971e8bc73876dc25ef40a575a12e91e4c033d8', false, '2026-08-11 03:35:48.373+00', '2026-07-11 20:35:48.373551+00'),
	(179, 27, '91648a08685ea6a0def54395ae64fb8907490d055090840143f4b661b7b78126', true, '2026-08-12 10:42:37.602578+00', '2026-07-12 20:42:37.603725+00'),
	(180, 27, '1be447d0245e16694d1413fbff068421b66aa9aaf2ab7e59e40a036d9f43457a', true, '2026-08-12 10:43:29.853916+00', '2026-07-12 20:43:29.854235+00'),
	(181, 27, 'a0decbe7a606572377b546517adc1eba10cd74d8f38ea88a19e89aa639b6de07', true, '2026-08-12 10:44:58.68785+00', '2026-07-12 20:44:58.687883+00'),
	(182, 27, 'c731f4c258de910cb75634b0dc0de388da596bac96cfbe3bac3ca0b1d6859ab6', true, '2026-08-12 10:44:59.659334+00', '2026-07-12 20:44:59.659612+00'),
	(183, 27, '8d55abb433bb8db82f574735900ea28332242c1261028d8cdab03a85c5ae2aed', true, '2026-08-12 10:45:00.247218+00', '2026-07-12 20:45:00.24728+00'),
	(185, 27, '4948111a2647335e091452eeb0cf4853c7539a26df0e902b84ebb41d893c2e08', true, '2026-08-12 10:45:05.274159+00', '2026-07-12 20:45:05.274265+00'),
	(186, 27, 'ea6db47efb44eb26aacf421cf7349b63b01d9a2159db6f6b476317ad7df38886', true, '2026-08-12 10:45:07.379191+00', '2026-07-12 20:45:07.379228+00'),
	(188, 27, '13161fa1c548e39b0b25f8f85ca61737050a61dc00a63c31b368a3a48331c340', true, '2026-08-12 10:45:10.554371+00', '2026-07-12 20:45:10.55456+00'),
	(189, 27, 'b03f26a9baec4228b5a37567faae3a1bdfe045677401a0f0c275f6986c6695c9', true, '2026-08-12 11:12:32.705502+00', '2026-07-12 21:12:32.705839+00'),
	(190, 27, '8b47dc2f2e39715704db0102c1b95454fec057ad8313407e166a592009185327', true, '2026-08-12 11:34:38.745946+00', '2026-07-12 21:34:38.745998+00'),
	(191, 27, '0b7b768f9094ffa32bd0331717b2f865c21cbfcb78ea8bf0c1ed88c539943cee', true, '2026-08-12 11:49:27.318478+00', '2026-07-12 21:49:27.318623+00'),
	(192, 27, 'ef26c1e6d5aded1f310b07014907a8fadd9ec0669d8f53c5faf346ac1bc283a5', true, '2026-08-12 12:00:40.505197+00', '2026-07-12 22:00:40.505441+00'),
	(193, 27, '2652da101e7718ebae96bd94fe89440b9244ea7026f0a2b5ed29a8163f0327c7', true, '2026-08-12 12:21:55.189422+00', '2026-07-12 22:21:55.189555+00'),
	(194, 27, '44e59829ac3e8b3b6e07a1bfff942ceb13210117a77cee75ceb632a9735a945e', true, '2026-08-12 12:34:14.717606+00', '2026-07-12 22:34:14.717637+00'),
	(195, 27, '854ca56320d4ccbee383c62051a723ff16d363be83c02e53a861a7f585aadb83', true, '2026-08-12 13:48:59.554327+00', '2026-07-12 23:48:59.554513+00'),
	(196, 27, '41438e5667e2163d8698633181539150886b7f7cdd74316ee8a4b7869d0153f2', true, '2026-08-12 14:24:14.005719+00', '2026-07-13 00:24:14.006059+00'),
	(197, 27, '204bec93de13f640945836e953c583919d71398feed9691207baf51e3f040813', true, '2026-08-12 14:24:24.362884+00', '2026-07-13 00:24:24.363247+00'),
	(198, 27, '2646ace1c2c5317fb6430e608eac18d3281eb1776798b3d77a2c442603caf0c9', true, '2026-08-12 14:29:46.947512+00', '2026-07-13 00:29:46.947834+00'),
	(199, 27, 'f08807dd1c461423672929612b2785a9e7fc1f1c7f760e6a4edaa4d82ba2f8f2', true, '2026-08-12 14:29:56.960435+00', '2026-07-13 00:29:56.960573+00'),
	(200, 27, '211f29328a35d6f2211d0fefac9133f7e5df1178687998f5f8fe1239a3a72421', true, '2026-08-12 14:37:57.389722+00', '2026-07-13 00:37:57.389948+00'),
	(201, 27, '2400f12ba34fcb8a3a07f011b59e2803591cb595268d4281dfc05497779a6482', true, '2026-08-12 14:38:22.38454+00', '2026-07-13 00:38:22.384937+00'),
	(202, 27, '1de34b13bedd91c3f2944d7781719a1f8714b9a7f8565c2192dbc7e5591943a6', true, '2026-08-12 14:45:03.802555+00', '2026-07-13 00:45:03.802847+00'),
	(203, 27, 'a83bafe90cd0830354a74dfc20bd323ded54e05e482863341adb2f0ab855a6a8', true, '2026-08-12 14:45:28.620412+00', '2026-07-13 00:45:28.620533+00'),
	(204, 27, '2fc077d297452dca04a60c40113861252bd58af546eb6412b93496df52db2d6e', false, '2026-08-12 14:46:46.744133+00', '2026-07-13 00:46:46.744173+00'),
	(206, 27, 'dc88158bbd5a7dc0652f1eed9787dff27d9e3a0158eeec3d50578ba564e0e144', true, '2026-08-12 14:46:53.741681+00', '2026-07-13 00:46:53.741818+00'),
	(207, 27, '063c654b0a98c9c7c8c8d17800ad4f204c3b29584f648498f7963334583b9268', false, '2026-08-12 14:46:58.358596+00', '2026-07-13 00:46:58.358628+00'),
	(209, 27, 'c139861bac00ca7bdee10ab22dc9bd5637293adbd8e3edbef77717eb8f42c3ce', true, '2026-08-12 14:47:05.452774+00', '2026-07-13 00:47:05.453462+00'),
	(210, 27, 'e8d8cae1daf0dd27024cd6bef8a843b3ec92ceb97fdb7cf1cfab77dda8fc5caf', false, '2026-08-12 14:47:14.612002+00', '2026-07-13 00:47:14.612031+00'),
	(212, 27, 'b21303ffd8abe387f3faba812dfca1021cc31c661fd9e899024a7bf7036d199c', true, '2026-08-12 14:48:30.679735+00', '2026-07-13 00:48:30.68072+00'),
	(213, 27, '79a0c2c70202f0d027672d794e51aae931e60566330b0d3af0904925b690e67e', true, '2026-08-12 14:48:32.619821+00', '2026-07-13 00:48:32.619849+00'),
	(214, 27, 'e131cbb6a338f4f1a734c0f34fa0faf3876d0bf356c0ffa3f742aa250319c62c', true, '2026-08-12 14:48:43.337927+00', '2026-07-13 00:48:43.337968+00'),
	(215, 27, 'db8a71c77b6b2a766d0fefefe3a62e3b80b2719954c9f4cbffc9d9036b569011', true, '2026-08-12 15:18:27.888405+00', '2026-07-13 01:18:27.890325+00'),
	(216, 27, '5ebd283c07b9e8d350d1460434aa286fbd05a848e223b7af8dc0d7301f5d1b2f', true, '2026-08-12 17:45:45.607707+00', '2026-07-13 03:45:45.608148+00'),
	(217, 27, '6cd5efd2ceb5678497a1075256ca20204a5da49ce885d2cec559c08c2b32b262', true, '2026-08-12 18:40:16.832812+00', '2026-07-13 04:40:16.833006+00'),
	(218, 27, '933d4774b5448538be366e0461bf8b21c4a0a56e3e7f1a2bc992e23310d955a1', true, '2026-08-12 19:47:10.456278+00', '2026-07-13 05:47:10.457493+00'),
	(219, 27, '8aab67982f36eba8c5da6b53fb6f95b86a21fbaf43c3074307f95b3e2c7218db', true, '2026-08-13 05:56:36.026508+00', '2026-07-13 15:56:36.026656+00'),
	(220, 27, 'e7aa6ad59f474c3d2afcbdce2e8c95d1329f30b05651e36b06c2e06237f42ab7', true, '2026-08-13 07:02:14.623772+00', '2026-07-13 17:02:14.623904+00'),
	(221, 27, 'ab6070eeee0a6146128936f879cf1646ba833a095b3c3ed1e4217674d1ae7b4b', true, '2026-08-13 07:43:34.148948+00', '2026-07-13 17:43:34.149127+00'),
	(222, 27, '597ebb5c880296d4f6fe9ed7c43cbf4b242504395b6243579f697a1943c18159', true, '2026-08-13 08:04:30.66263+00', '2026-07-13 18:04:30.662758+00'),
	(223, 27, 'e2158544d65c48802c9d1db5ed44f4241009af3e729dc11b6aa29c78c8d74fe1', true, '2026-08-13 08:40:49.676256+00', '2026-07-13 18:40:49.676428+00'),
	(224, 27, '1a2e81610f1986b193a5ecaf0e82734590e324c86a1591dd9ac0aec8a3ea5a89', true, '2026-08-13 11:40:21.927854+00', '2026-07-13 21:40:21.928341+00'),
	(225, 27, '06f7b8daa314f97b506cf482451c7cd9f973a19995f513a9553c0bbaa44b486c', true, '2026-08-13 12:27:05.431333+00', '2026-07-13 22:27:05.431582+00'),
	(226, 27, '0c1893eb1f635d240ba09a4ed00b4c7ea4902aa14dbc1aae3bdb5e857f63c034', true, '2026-08-13 13:40:44.344615+00', '2026-07-13 23:40:44.344793+00'),
	(227, 27, '7cfba424f64fdbb9cf58b426a6f247058aa5f068ba6afd765f5bea90cb665e17', true, '2026-08-13 13:51:20.857175+00', '2026-07-13 23:51:20.85729+00'),
	(228, 27, '77a713c0c64f89b7d0d55a01aaac13bc976dbf5aeb27505f0395856b3dab15fb', true, '2026-08-13 16:42:18.954541+00', '2026-07-14 02:42:18.955891+00'),
	(229, 27, '206efb2f511d55f0a98615fec3ae137f1b65e06738cb446aa8a63c319747fd4b', true, '2026-08-13 16:48:27.092107+00', '2026-07-14 02:48:27.09237+00'),
	(230, 27, '11c6168adbea84089da567ef77756f3a84759eeb57ac4ab96a5818b46da38a83', true, '2026-08-14 07:59:55.363432+00', '2026-07-14 17:59:55.365192+00'),
	(231, 27, 'c475bddbce03147112746144caf93f354ac57398dc6f52291321cfc531ac94e0', true, '2026-08-15 12:31:43.384757+00', '2026-07-15 22:31:43.384986+00'),
	(232, 27, 'd465c47a51fc7d047e0b49f5d42bd73a1903ff560ca6cac810e2c75cb18110b9', true, '2026-08-15 14:24:44.617904+00', '2026-07-16 00:24:44.618074+00'),
	(233, 27, '41bcb62d1963de174a94fa3d8dc992d5e60749116d39420cbcf5f6fe2a92e8ad', true, '2026-08-15 15:41:59.284361+00', '2026-07-16 01:41:59.28453+00'),
	(234, 27, 'e69e9c34bad2f769f3b63da5e605664a2b7713e220de04e1180c43e01c4a3445', true, '2026-08-15 15:50:49.3203+00', '2026-07-16 01:50:49.320407+00'),
	(235, 27, '3a83a349e67f0cc661e3fcd61a316ba75c0a37d90a1fae50eac8f6b0ff01b674', true, '2026-08-15 16:21:41.864108+00', '2026-07-16 02:21:41.864205+00'),
	(236, 27, 'cea3425439c7653282c1686ef26acfd3b61ccc0a259c25b2e86bb183c39df936', true, '2026-08-15 16:51:05.563786+00', '2026-07-16 02:51:05.563883+00'),
	(237, 27, 'e533e6ce0bfa43e7e6becf8576636cdd36198e10c718b6b9be3a981c6598aed8', true, '2026-08-15 19:50:17.710687+00', '2026-07-16 05:50:17.712743+00'),
	(238, 27, '20f26c7548f8282f8777a0eeb5e604db4f5f0fa40fe4fb01d7028f4a5a3b99ff', true, '2026-08-15 20:57:22.449282+00', '2026-07-16 06:57:22.449551+00'),
	(239, 27, '22460cc51e73373a26650c8704753092056df4fbc78b3bb1499193305e8e597c', true, '2026-08-16 09:24:21.780342+00', '2026-07-16 19:24:21.780476+00'),
	(240, 27, 'd8c97d70f6463e5904f1d5236aa04cd079e2e684fbe940f4474eecf9900d54ef', true, '2026-08-16 10:33:02.361029+00', '2026-07-16 20:33:02.361165+00'),
	(241, 27, 'f4922d3147ae9b29d311f955ddafb675480d462bdb9d173b508a3a79a81f0406', true, '2026-08-16 11:01:06.159895+00', '2026-07-16 21:01:06.160058+00'),
	(242, 27, 'b2fbf684b73f75607ba006c16c91646745d58801019a10b57dd3c37377b74ed9', true, '2026-08-16 11:30:41.69982+00', '2026-07-16 21:30:41.699944+00'),
	(243, 27, '0a369553251b42721f9b6b723bc948d49c009c03e387b1e316936dcdba1b70e6', true, '2026-08-16 13:02:46.749408+00', '2026-07-16 23:02:46.74965+00'),
	(244, 27, 'ff560f8fa96b77c875bf4e77123d6866b1319072f1138718dba8e4f8e4082750', true, '2026-08-16 15:26:47.036451+00', '2026-07-17 01:26:47.036591+00'),
	(245, 27, '76e88067a7da25323e280412d5d100097aa90f16206e1851f31b011251b6af10', true, '2026-08-16 16:20:25.315557+00', '2026-07-17 02:20:25.315823+00'),
	(246, 27, '5e8340ad399a5dd8eb1246f73186c86ca7a6a7156f211bbd06e2a88ed736e867', true, '2026-08-16 16:39:39.48316+00', '2026-07-17 02:39:39.483363+00'),
	(247, 27, '1fd3cfa55fa3ed97ad04381fc4dd08ada6284d75a9c166a8a45f469278bef19e', true, '2026-08-17 07:01:32.842629+00', '2026-07-17 17:01:32.842784+00'),
	(248, 27, 'f712532e0314162b5a0984960ec15d26f387bf4f1b6b7a15aa8d6284b5cce6f9', true, '2026-08-17 09:34:42.413098+00', '2026-07-17 19:34:42.413259+00'),
	(249, 27, 'bb9442a2cbdab92469e607a4d327b97c68b754a325a0dc3a274d53e164159a5f', true, '2026-08-17 10:12:19.216612+00', '2026-07-17 20:12:19.216777+00'),
	(250, 27, '1c52e62a9eea1f48b37bb6f987aefb25a3e321143de126648831924ad4cf4328', true, '2026-08-18 08:04:48.626489+00', '2026-07-18 18:04:48.62664+00'),
	(251, 27, 'd292b67bf55ff9530df57a1161847822a86ecd02cc0cc6c3e4a82bf488ec2d62', true, '2026-08-18 14:17:47.608172+00', '2026-07-19 00:17:47.608351+00'),
	(252, 27, 'e6790e10e48c833d5c1dae7b26e45ba2acf33f38997a9f7b2e26f7753b87b8b6', true, '2026-08-19 08:52:11.684321+00', '2026-07-19 18:52:11.684495+00'),
	(253, 27, 'eac0123c1658f1a7e45bf1afeb0ded62352f685526b49018dd5502afb62d2942', true, '2026-08-19 09:33:08.334451+00', '2026-07-19 19:33:08.334624+00'),
	(254, 27, '48a9625f18a5f61e462ca6202e164d54cdf8c58d8ec7812a6d55b7510571c63f', true, '2026-08-19 10:45:19.190523+00', '2026-07-19 20:45:19.190649+00'),
	(255, 27, 'a37738828fd553fa598da4a37d8c8fa82d2cc8e4cdf24de3512a3da5df9f473c', true, '2026-08-19 10:55:29.795282+00', '2026-07-19 20:55:29.795387+00'),
	(256, 27, '40d5b44a7f9f08bee6f70dcf8eb17873e4f997ae4d8907f728a2a1602b7b0cc4', true, '2026-08-19 12:21:40.04593+00', '2026-07-19 22:21:40.046152+00'),
	(257, 27, 'f3272ac95797567acd4a9c510dd2baef67141b9932b5e44ba6190111de554e98', true, '2026-08-19 12:39:23.972832+00', '2026-07-19 22:39:23.973001+00'),
	(258, 27, '4b5122fc3c5a9da0a7df2dab74dd2389305f70efa424db191f0ee1e64b834e2e', true, '2026-08-19 14:51:28.075084+00', '2026-07-20 00:51:28.075345+00'),
	(259, 27, '860680df962b7ebaf08bf2d708cfb3f7dd90772e07ef98136c007658095c5073', true, '2026-08-19 15:42:11.510962+00', '2026-07-20 01:42:11.51113+00'),
	(260, 27, '1807ed81055df9ce95338a4b498b05dea17b1d72d55dfe33529a325a2f91a47d', true, '2026-08-19 17:34:06.384203+00', '2026-07-20 03:34:06.384367+00'),
	(261, 27, 'd9e7e8c1454b0f276bd05db51ea9db46fc118bbc45729ad88eb1ee6c3eb57134', true, '2026-08-19 20:43:14.960463+00', '2026-07-20 06:43:14.961065+00'),
	(262, 27, 'e6b5db535ba9e62028f860b3e6fc5fbc2e85ab4efe0f97f42c9ee20ed4b78f82', true, '2026-08-20 09:47:19.966544+00', '2026-07-20 19:47:19.966671+00'),
	(263, 27, '3c8196dd21dd9b4e351cd8f9e03f125ac19808f3f85d4c413f5e7c22ad161a46', true, '2026-08-20 09:54:13.368752+00', '2026-07-20 19:54:13.368972+00'),
	(264, 27, '811bd1c1e8e698b07e5c3d84efb3e192fb5fd999b5fc03bef94f0d50ca939917', true, '2026-08-20 11:05:57.969277+00', '2026-07-20 21:05:58.011511+00'),
	(265, 27, 'df74301bdd6a2c75eda8c3a9d0396b65e0c15467e5fa7c1069a2d4986b39ec13', false, '2026-08-20 13:37:52.07138+00', '2026-07-20 23:37:52.073285+00')
ON CONFLICT (token) DO UPDATE SET
	user_id = EXCLUDED.user_id,
	is_revoked = EXCLUDED.is_revoked,
	expires_at = EXCLUDED.expires_at,
	created_at = EXCLUDED.created_at;

SELECT setval(pg_get_serial_sequence('users', 'id'), COALESCE((SELECT MAX(id) FROM users), 1), true);
SELECT setval(pg_get_serial_sequence('teams', 'id'), COALESCE((SELECT MAX(id) FROM teams), 1), true);
SELECT setval(pg_get_serial_sequence('team_members', 'id'), COALESCE((SELECT MAX(id) FROM team_members), 1), true);
SELECT setval(pg_get_serial_sequence('projects', 'id'), COALESCE((SELECT MAX(id) FROM projects), 1), true);
SELECT setval(pg_get_serial_sequence('project_teams', 'id'), COALESCE((SELECT MAX(id) FROM project_teams), 1), true);
SELECT setval(pg_get_serial_sequence('tasks', 'id'), COALESCE((SELECT MAX(id) FROM tasks), 1), true);
SELECT setval(pg_get_serial_sequence('task_assignments', 'id'), COALESCE((SELECT MAX(id) FROM task_assignments), 1), true);
SELECT setval(pg_get_serial_sequence('comments', 'id'), COALESCE((SELECT MAX(id) FROM comments), 1), true);
SELECT setval(pg_get_serial_sequence('attachments', 'id'), COALESCE((SELECT MAX(id) FROM attachments), 1), true);
SELECT setval(pg_get_serial_sequence('refresh_tokens', 'id'), COALESCE((SELECT MAX(id) FROM refresh_tokens), 1), true);

COMMIT;
