using Ardalis.Specification;
using eshopPractice.ApplicationCore.Entities.BasketAggregate;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;

namespace eshopPractice.ApplicationCore.Specifications
{
    public sealed class BasketWithItemsSpecification : Specification<Basket>
    {
        public BasketWithItemsSpecification(int basketId)
        {
            Query
                .Where(b => b.Id == basketId)
                .Include(b => b.Items);
        }

        public BasketWithItemsSpecification(string buyerId)
        {
            Query
                .Where(b => b.BuyerId == buyerId)
                .Include(b => b.Items);
        }
    }
}
