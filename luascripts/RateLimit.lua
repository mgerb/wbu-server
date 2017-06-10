-- key takes in <key>:<ip address>
local key = KEYS[1]

local rateCount = redis.call("INCR", key)

-- rateCount will be one if key previously did not exist
if rateCount == 1 then
    redis.call("EXPIRE", key, 60)
end

-- currently set the rate limit at 200 requests per minute
if rateCount > 200 then
    return redis.error_reply("Rate exceeded")
end

return "Success"
