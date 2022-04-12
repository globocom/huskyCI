using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Linq;
using System.Security.Policy;
using System.Threading.Tasks;
using eshopPractice.ApplicationCore.Interfaces;
using eshopPractice.Interfaces;
using eshopPractice.ViewModels;
using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.Mvc.RazorPages;

namespace eshopPractice.Pages.Basket
{
    public class IndexModel : PageModel
    {
        private readonly IBasketService _basketService;
        private readonly IBasketViewModelService _basketViewModelService;
        private string _username = "tempuser";
        public IndexModel(IBasketService basketService, IBasketViewModelService basketViewModelService)
        {
            _basketService = basketService;
            _basketViewModelService = basketViewModelService;
        }

        public BasketViewModel BasketModel { get; set; } = new BasketViewModel();

        public async Task OnGet()
        {
            await SetBasketModelAsync();
        }

        public async Task<IActionResult> OnPost(CatalogItemViewModel productDetails)
        {
            if (productDetails?.Id == null)
            {
                return RedirectToPage("/Index");
            }
            await SetBasketModelAsync();

            await _basketService.AddItemToBasket(BasketModel.Id, productDetails.Id, productDetails.Price);

            await SetBasketModelAsync();

            return RedirectToPage();
        }

        public async Task OnPostUpdate(IEnumerable<BasketItemViewModel> items)
        {
            await SetBasketModelAsync();

            if (!ModelState.IsValid)
            {
                return;
            }

            var updateModel = items.ToDictionary(b => b.Id.ToString(), b => b.Quantity);
            await _basketService.SetQuantities(BasketModel.Id, updateModel);

            await SetBasketModelAsync();
        }

        private async Task SetBasketModelAsync()
        {
            GetOrSetBasketCookieAndUserName();
            BasketModel = await _basketViewModelService.GetOrCreateBasketForUser(_username);
        }

        private void GetOrSetBasketCookieAndUserName()
        {
            if (Request.Cookies.ContainsKey(Constants.BASKET_COOKIENAME))
            {
                _username = Request.Cookies[Constants.BASKET_COOKIENAME];
            }

            if (_username != null) return;

            _username = Guid.NewGuid().ToString();
            var cookieOptions = new CookieOptions { IsEssential = true };
            cookieOptions.Expires = DateTime.Today.AddYears(10);
            Response.Cookies.Append(Constants.BASKET_COOKIENAME, _username, cookieOptions);
        }
        
    }
}
