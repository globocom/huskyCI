using eshopPractice.Interfaces;
using eshopPractice.Pages.Basket;
using eshopPractice.ViewModels;
using Microsoft.AspNetCore.Mvc;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;

namespace eshopPractice.Pages.Shared.Components.Basket
{
    public class Basket : ViewComponent
    {
        private readonly IBasketViewModelService _basketService;
        public Basket(IBasketViewModelService basketService)
        {
            _basketService = basketService;
        }

        public async Task<IViewComponentResult> InvokeAsync(string userName)
        {
            var vm = new BasketComponentViewModel();
            vm.ItemsCount = (await GetBasketViewModelAsync()).Items.Sum(i => i.Quantity);
            return View(vm);
        }

        private async Task<BasketViewModel> GetBasketViewModelAsync()
        {
            string anonymousId = GetBasketIdFromCookie();
            if (anonymousId == null) return new BasketViewModel();
            return await _basketService.GetOrCreateBasketForUser(anonymousId);
        }

        private string GetBasketIdFromCookie()
        {
            if (Request.Cookies.ContainsKey(Constants.BASKET_COOKIENAME))
            {
                return Request.Cookies[Constants.BASKET_COOKIENAME];
            }
            return null;
        }
    }
}
