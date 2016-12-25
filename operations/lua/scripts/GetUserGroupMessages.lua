local key = KEYS[1]

local messages = redis.call("SMEMBERS", key)

-- remove the messages after retreiving
redis.call("DEL", key)

return messages