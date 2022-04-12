using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using eshopPractice.Interfaces;
using eshopPractice.ViewModels;
using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.Mvc.RazorPages;

namespace eshopPractice.Pages
{
    public class IndexModel : PageModel
    {
        private readonly ICatalogViewModelService _catalogViewModelService;

        public IndexModel(ICatalogViewModelService catalogViewModelService)
        {
            _catalogViewModelService = catalogViewModelService;
        }

        public CatalogIndexViewModel CatalogModel { get; set; } = new CatalogIndexViewModel(); 
        public async Task OnGet(CatalogIndexViewModel catalogModel, int? pageId)
        {
            CatalogModel = await _catalogViewModelService.GetCatalogItems(pageId ?? 0, Constants.ITEMS_PER_PAGE, catalogModel.BrandFilterApplied, catalogModel.TypesFilterApplied);
        }
    }
}
