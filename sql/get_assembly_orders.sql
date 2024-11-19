CREATE DEFINER=`pmp`@`%` PROCEDURE `get_assembly_orders`(IN dt_start datetime, IN dt_finish datetime)
    COMMENT 'Получение списка собранных заказов за выбранный период'
BEGIN
    SELECT 
        `o`.`order_uid` AS `order_uid`,
        `o`.`order_date` AS `order_date`,
        `o`.`order_sum` AS `order_sum`,
        `o`.`folio_num` AS `folio_num`,
        `o`.`unicum_num` AS `unicum_num`,
        `o`.`folio_date` AS `folio_date`,
        `o`.`user_id` AS `user_id`,
        `o`.`employee_id` AS `employee_id`,
        `o`.`created_at` AS `created_at`,
        `o`.`updated_at` AS `assembly_date`,
        TIMESTAMPDIFF(MINUTE,
            `o`.`created_at`,
            `o`.`updated_at`) AS `date_diff_minutes`,
        ROUND((TIMESTAMPDIFF(MINUTE,
                    `o`.`created_at`,
                    `o`.`updated_at`) / 60),
                2) AS `date_diff_hours`,
        CONCAT(`u`.`first_name`,
                ' ',
                `u`.`name`,
                ' ',
                `u`.`last_name`) AS `user_name`,
        CONCAT(`e`.`first_name`,
                ' ',
                `e`.`name`,
                ' ',
                `e`.`last_name`) AS `employee_name`,
        `o`.`client_name` AS `client_name`,
        `o`.`vid_doc` AS `vid_doc`
    FROM
        ((`orders` `o`
        LEFT JOIN `users` `u` ON ((`u`.`id` = `o`.`user_id`)))
        LEFT JOIN `employees` `e` ON ((`e`.`id` = `o`.`employee_id`)))
    WHERE
        (IFNULL(`o`.`user_id`, 0) <> 0)
        and CAST(`o`.`folio_date` as date) between dt_start and dt_finish
        and `o`.`deleted_at` IS NULL;	
END