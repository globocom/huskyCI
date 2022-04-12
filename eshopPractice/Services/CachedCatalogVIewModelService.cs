using eshopPractice.Extensions;
using eshopPractice.Interfaces;
using eshopPractice.ViewModels;
using Microsoft.AspNetCore.Mvc.Rendering;
using Microsoft.Extensions.Caching.Memory;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;

namespace eshopPractice.Services
{
    public class CachedCatalogVIewModelService : ICatalogViewModelService
    {
        private readonly IMemoryCache _cache;
        private readonly ICatalogViewModelService _catalogViewModelService;

        public CachedCatalogVIewModelService(IMemoryCache cache, CatalogViewModelService catalogViewModelService)
        {
            _cache = cache;
            _catalogViewModelService = catalogViewModelService;
        }

        public async Task<IEnumerable<SelectListItem>> GetBrands()
        {
            return await _cache.GetOrCreateAsync(CacheHelpers.GenerateBrandsCacheKey(), async entry =>
            {
                entry.SlidingExpiration = CacheHelpers.DefaultCacheDuration;
                return await _catalogViewModelService.GetBrands();
            });
        }

        public async Task<CatalogIndexViewModel> GetCatalogItems(int pageIndex, int itemsPage, int? brandId, int? typeId)
        {
            var cacheKey = CacheHelpers.GenerateCatalogItemCacheKey(pageIndex, Constants.ITEMS_PER_PAGE, brandId, typeId);

            return await _cache.GetOrCreateAsync(cacheKey, async entry =>
            {
                entry.SlidingExpiration = CacheHelpers.DefaultCacheDuration;
                return await _catalogViewModelService.GetCatalogItems(pageIndex, itemsPage, brandId, typeId);
            });
        }

        public async Task<IEnumerable<SelectListItem>> GetTypes()
        {
            return await _cache.GetOrCreateAsync(CacheHelpers.GenerateTypesCacheKey(), async entry => 
            {
                entry.SlidingExpiration = CacheHelpers.DefaultCacheDuration;
                return await _catalogViewModelService.GetTypes();
            });
        }
    }
}
