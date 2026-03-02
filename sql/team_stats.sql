select t.id, t.title, COUNT(DISTINCT tm.user_id) as member_count, COUNT(DISTINCT task.id) AS done_tasks_last_7_days
from tasks t
         left join team_members tm on tm.team_id = t.id
         left join tasks task
                   on task.team_id = t.id and task.status = 'done' and task.created_at >= now() - interval 7 day
group by t.id, t.title;



select t.id, t.team_id, t.assignee_id
from tasks t
         left join team_members tm on tm.team_id = t.team_id
    and tm.user_id = t.assignee_id
where tm.user_id is null
  and t.assignee_id is not null