CREATE TABLE execution_trackers (
    id bigint(11) UNSIGNED NOT NULL AUTO_INCREMENT,
    name varchar(255) NOT NULL,
    last_executor varchar(255),
    next_execution_time timestamp(6),
    created_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
) ENGINE=InnoDB;

CREATE UNIQUE INDEX uqx_name ON execution_trackers (name);
