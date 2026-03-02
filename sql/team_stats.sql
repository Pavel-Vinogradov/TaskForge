SELECT t.id                       AS team_id,
       t.name                     AS team_name,
       COUNT(DISTINCT tm.user_id) AS member_count,
       COUNT(DISTINCT CASE
                          WHEN task.status = 'done'
                              AND task.created_at >= '2026-03-01 00:00:00'
                              THEN task.id
           END)                   AS completed_tasks_count,
       '2026-03-01 00:00:00'      AS period_start,
       '2026-03-31 23:59:59'      AS period_end
FROM teams t
         LEFT JOIN team_members tm ON t.id = tm.team_id
         LEFT JOIN tasks task ON t.id = task.team_id
GROUP BY t.id, t.name
ORDER BY t.name;