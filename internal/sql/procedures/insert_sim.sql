DELIMITER //

CREATE PROCEDURE AddSim(
	IN simNumber VARCHAR(15),
    IN providerId INT,
    IN isActivated TINYINT,
    IN activateUntil BIGINT,
    IN isBlocked TINYINT,
    OUT lastId INT
)
BEGIN
    -- Start a new transaction
    START TRANSACTION;

    -- Insert the record
    INSERT INTO sim (number, provider_id, is_activated, activate_until, is_blocked)
    VALUES (simNumber, providerId, isActivated, activateUntil, isBlocked);

    -- Get the last inserted ID
    SET lastId = LAST_INSERT_ID();

    -- Check if the insert was successful
    IF lastId > 0 THEN
        -- Commit the transaction
        COMMIT;
    ELSE
        -- Rollback the transaction
        ROLLBACK;
        -- Raise an error
    END IF;
END //

DELIMITER ;