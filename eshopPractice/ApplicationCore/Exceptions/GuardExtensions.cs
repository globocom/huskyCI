using eshopPractice.ApplicationCore.Entities.BasketAggregate;
using eshopPractice.ApplicationCore.Exceptions;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;

namespace Ardalis.GuardClauses
{
    public static class BasketGuards
    {
        public static void NullBasket(this IGuardClause guardClause, int basketId, Basket basket)
        {
            if (basket == null)
            {
                throw new BasketNotFoundException(basketId);
            }
        }

        public static void EmptyBasketOnCheckout(this IGuardClause guardClause, IReadOnlyCollection<BasketItem> basketItems)
        {
            if (!basketItems.Any())
                throw new EmptyBasketOnCheckoutException();
        }
    }
}
