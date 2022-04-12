using Ardalis.Specification;
using eshopPractice.ApplicationCore.Entities;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;

namespace eshopPractice.ApplicationCore.Specifications
{
    public class CatalogItemsSpecification : Specification<CatalogItem>
    {
        public CatalogItemsSpecification(params int[] ids)
        {
            Query.Where(c => ids.Contains(c.Id));
        }
    }
}
