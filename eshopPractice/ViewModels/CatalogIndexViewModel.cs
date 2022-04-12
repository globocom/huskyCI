using Microsoft.AspNetCore.Mvc.Rendering;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;

namespace eshopPractice.ViewModels
{
    public class CatalogIndexViewModel
    {
        public List<CatalogItemViewModel> CatalogItems { get; set; }
        public List<SelectListItem> Brands { get; set; }
        public List<SelectListItem> Types { get; set; }
        public int? BrandFilterApplied { get; set; }
        public int? TypesFilterApplied { get; set; }
        public PaginationInfoViewModel PaginationInfo { get; set; }
    }
}
