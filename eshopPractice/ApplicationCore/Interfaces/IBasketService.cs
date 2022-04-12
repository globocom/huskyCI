using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;

namespace eshopPractice.ApplicationCore.Interfaces
{
    public interface IBasketService
    {
        Task TransferBasketAsync(string anonymousId, string userName);
        Task AddItemToBasket(int basketId, int catalogItemId, decimal price, int quantity = 1);
        Task SetQuantities(int basketId, Dictionary<string, int> quantities);
        Task DeleteBasketAsync(int basketId);
    }
}
