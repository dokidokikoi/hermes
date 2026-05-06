-- Example Lua Scraper Plugin
-- This is a template for creating new scraper plugins using Lua

-- Plugin metadata
local name = "example_scraper"
local version = "1.0.0"

-- Return plugin information
function get_info()
	return {
		name = name,
		version = version
	}
end

-- Search for games by keyword
-- @param keyword: The search keyword
-- @param page: The page number (1-indexed)
-- @return: Array of search results, each containing:
--   - name: The game name
--   - key: Unique identifier
--   - url: The URL to the game page
--   - summary: Brief description (optional)
--   - cover: Cover image URL (optional)
function search(keyword, page)
	log.info("Searching for: " .. keyword .. " on page " .. tostring(page))

	-- Example: Make a request to a search endpoint
	local status, body = http.get("https://example.com/search?q=" .. keyword .. "&page=" .. tostring(page), scraper_headers, scraper_proxy)

	if status ~= 200 then
		log.error("Search failed with status: " .. tostring(status))
		return {}
	end

	-- Parse the HTML response
	local doc = html.parse(body)
	local results = {}

	-- Extract search results from the page
	doc:select(".game-item"):each(function(i, elem)
		local result = {}

		-- Extract game name
		local nameElem = elem:select(".title")
		if nameElem then
			result.name = nameElem:text()
		end

		-- Extract URL
		local linkElem = elem:select("a")
		if linkElem then
			result.url = linkElem:attr("href")
			result.key = result.url -- Use URL as key
		end

		-- Extract summary
		local summaryElem = elem:select(".summary")
		if summaryElem then
			result.summary = summaryElem:text()
		end

		-- Extract cover image
		local imgElem = elem:select("img")
		if imgElem then
			result.cover = imgElem:attr("src")
		end

		table.insert(results, result)
	end)

	return results
end

-- Get detailed game information from a URL
-- @param url: The URL to the game page
-- @return: Game item table containing:
--   - vndb_id: VNDB ID (optional)
--   - jan_code: JAN code (optional)
--   - dl_code: DL code (optional)
--   - name: Game name
--   - alias: Array of alternate names (optional)
--   - cover: Cover image URL
--   - images: Array of screenshot URLs (optional)
--   - category: { name = category_name }
--   - series: Array of { name = series_name }
--   - brands: Array of { vndb_id = "", name = "", links = {} }
--   - price: Price string (optional)
--   - issue_date: Issue date in ISO format (optional)
--   - story: Game description
--   - platform: Platform string (optional)
--   - tags: Array of { name = "", lang = "" }
--   - characters: Array of character tables
--   - staff: Array of staff tables
--   - links: Array of { name = "", url = "", type = "" }
--   - other_info: Additional info (optional)
function get_item(url)
	log.info("Getting game details from: " .. url)

	-- Make a request to the game page
	local status, body = http.get(url, scraper_headers, scraper_proxy)

	if status ~= 200 then
		log.error("Get item failed with status: " .. tostring(status))
		return nil
	end

	-- Parse the HTML response
	local doc = html.parse(body)

	-- Create the game item
	local item = {}

	-- Extract basic information
	local titleElem = doc:select(".game-title")
	if titleElem then
		item.name = titleElem:text()
	end

	-- Extract cover
	local coverElem = doc:select(".cover-image")
	if coverElem then
		item.cover = coverElem:attr("src")
	end

	-- Extract story/description
	local storyElem = doc:select(".story")
	if storyElem then
		item.story = storyElem:text()
	end

	-- Extract price
	local priceElem = doc:select(".price")
	if priceElem then
		item.price = priceElem:text()
	end

	-- Extract brand
	local brandElem = doc:select(".brand")
	if brandElem then
		item.brands = {
			{
				name = brandElem:text()
			}
		}
	end

	-- Extract tags
	item.tags = {}
	doc:select(".tags .tag"):each(function(i, elem)
		table.insert(item.tags, {
			name = elem:text(),
			lang = "ja" -- or detect language
		})
	end)

	-- Extract links
	item.links = {
		{
			name = name,
			url = url,
			type = "info"
		}
	}

	return item
end