use eastwh;

select ut.id,
	   ut.team_id,
       t.name,
       u.id as user_id,
	   CONCAT(u.first_name, ' ', u.name, ' ', u.last_name) AS user_name,
	   CONCAT(e.first_name, ' ', e.name, ' ', e.last_name) AS employee_name
from user_teams ut
	left join teams t on t.id = ut.team_id
    left join users u on u.id = ut.user_id
    left join employee_teams et on et.team_id = t.id
    left join employees e on e.id = et.employee_id
where ut.user_id = 12;	