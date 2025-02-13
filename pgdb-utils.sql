-- Disable and enable triggers
SET session_replication_role = replica;
SET session_replication_role = DEFAULT;

-- Rebuild platform strings

---- Dry Run

SELECT COUNT(*) as rows_that_would_change
FROM game
LEFT JOIN (
    SELECT game_id,
           string_agg((SELECT primary_alias FROM platform WHERE id = p.platform_id), '; ') AS new_platforms_str
    FROM game_platforms_platform p
    GROUP BY game_id
) subquery ON game.id = subquery.game_id
WHERE (game.platforms_str IS DISTINCT FROM subquery.new_platforms_str
       OR (game.platforms_str IS NULL AND subquery.new_platforms_str IS NOT NULL)
       OR (game.platforms_str IS NOT NULL AND subquery.new_platforms_str IS NULL));

--- Live Run

UPDATE game
SET platforms_str = new_platforms_str,
reason = 'Rebuild Platforms String',
user_id = 810112564787675166
FROM (
    SELECT id,
           string_agg((SELECT primary_alias FROM platform WHERE id = p.platform_id), '; ') AS new_platforms_str
    FROM game
    LEFT JOIN game_platforms_platform p ON p.game_id = game.id
    GROUP BY game.id
) subquery
WHERE game.id = subquery.id
  AND (game.platforms_str IS DISTINCT FROM subquery.new_platforms_str
       OR (game.platforms_str IS NULL AND subquery.new_platforms_str IS NOT NULL)
       OR (game.platforms_str IS NOT NULL AND subquery.new_platforms_str IS NULL));

-- Rebuild tag strings

-- TODO: Make sure to only update necessary rows

UPDATE game
SET tags_str = coalesce(
    (
        SELECT string_agg(
                       (SELECT primary_alias FROM tag WHERE id = t.tag_id), '; '
                   )
        FROM game_tags_tag t
        WHERE t.game_id = game.id
    ), ''
) WHERE 1=1

-- Update app paths

---- Dry Run

SELECT COUNT(*) as rows_that_would_change
FROM game_data
WHERE application_path = 'app_path_here';

SELECT COUNT(*) as rows_that_would_change
FROM game
WHERE application_path = 'app_path_here';

---- Live Run

UPDATE game_data
SET application_path = 'new_app_path'
WHERE application_path = 'app_path_here';

UPDATE game
SET reason = 'Update Application Path',
    user_id = 810112564787675166
WHERE id IN (
    SELECT game_id FROM game_data
    WHERE application_path = 'new_app_path'
);

UPDATE game
SET application_path = 'new_app_path',
    reason = 'Update Application Path',
    user_id = 810112564787675166
WHERE application_path = 'app_path_here';

-- Update Launch Commands

---- Dry Run

SELECT COUNT(*) as rows_that_would_change
FROM game_data
WHERE application_path = 'app_path_here'
AND launch_command NOT LIKE 'prefix%'
AND launch_command LIKE 'http%';

SELECT COUNT(*) as rows_that_would_change
FROM game
WHERE application_path = 'app_path_here'
AND launch_command NOT LIKE 'prefix%'
AND launch_command LIKE 'http%';

---- Live Run

UPDATE game_data
SET launch_command = concat('prefix ', launch_command)
WHERE application_path = 'app_path_here';

UPDATE game
SET reason = 'Update Launch Command',
    user_id = 810112564787675166
WHERE id IN (
    SELECT game_id FROM game_data
    WHERE application_path = 'app_path_here'
    AND launch_command LIKE 'prefix%'
);

UPDATE game
SET launch_command = concat('prefix ', launch_command),
    reason = 'Update Launch Command',
    user_id = 810112564787675166
WHERE application_path = 'app_path_here'
AND launch_command LIKE 'prefix%';