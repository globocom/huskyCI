using Ardalis.Specification;
using eshopPractice.ApplicationCore.Entities;
using System;
using System.Collections.Generic;
using System.Globalization;
using System.Linq;
using System.Threading.Tasks;

namespace eshopPractice.ApplicationCore.Specifications
{
    public class CatalogFilterSpecification : Specification<CatalogItem>
    {
        public CatalogFilterSpecification(int? brandId, int? typeId)
        {
            Query.Where(i => (!brandId.HasValue || i.CatalogBrandId == brandId) &&
                (!typeId.HasValue || i.CatalogTypeId == typeId));
        }
    }
}
