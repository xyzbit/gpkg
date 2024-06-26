package limiter

const (
	tokenLimiterLuaScript = `
local needTokens  = tonumber(ARGV[1])
local maxTokens = tonumber(ARGV[2]) -- 最大令牌数量
local tokensPerSecond = tonumber(ARGV[3])  -- 每秒产生的令牌数量
local nowTs = tonumber(ARGV[4]) -- 当前时间
local storedTokens = 0
local lastTs = 0

local exist = redis.call('EXISTS', KEYS[1])
if exist == 0 then
    storedTokens = maxTokens
    lastTs = nowTs
else
    local limiterInfo = redis.call('HMGET', KEYS[1], 'last_ts', 'stored_tokens')
    lastTs = limiterInfo[1]
    storedTokens = limiterInfo[2] -- 当前剩余的令牌数量
end

local returnTokens = 0 -- 最终返回的令牌数量

local timePassTs = math.max(nowTs - lastTs, 0) -- 从上次获取令牌到现在经过的时间
storedTokens = math.min(storedTokens+(timePassTs*tokensPerSecond), maxTokens) -- 当前剩余的令牌数量，最大不能超过规定的数量

local newTokens = storedTokens
if storedTokens >= needTokens then
    returnTokens = needTokens
    newTokens = storedTokens - needTokens
end

-- 更新缓存
redis.call('HMSET', KEYS[1], 'last_ts', nowTs, 'stored_tokens', newTokens)
redis.call('EXPIRE', KEYS[1], (maxTokens/tokensPerSecond)*2)

return returnTokens
`
)
