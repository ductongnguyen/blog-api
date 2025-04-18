-- Disable foreign key checks to allow dropping tables in any order if needed,
-- but it's better practice to drop in reverse order of dependency.
SET FOREIGN_KEY_CHECKS=0;

-- Drop tables in reverse order of creation (or dependency)
DROP TABLE IF EXISTS `comments`;
DROP TABLE IF EXISTS `post_tags`;
DROP TABLE IF EXISTS `posts`;
DROP TABLE IF EXISTS `tags`;
DROP TABLE IF EXISTS `categories`;
DROP TABLE IF EXISTS `users`;

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS=1;