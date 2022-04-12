using Ardalis.Specification;
using eshopPractice.ApplicationCore.Entities;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;

namespace eshopPractice.ApplicationCore.Specifications
{
    public class CatalogFilterPaginatedSpecification : Specification<CatalogItem>
    {
        public CatalogFilterPaginatedSpecification(int skip, int take, int? brandId, int? typeId) : base()
        {
            Query
                .Where(i => (!brandId.HasValue || i.CatalogBrandId == brandId) &&
                    (!typeId.HasValue || i.CatalogTypeId == typeId))
                    .Paginate(skip, take);
        }
    }
}
