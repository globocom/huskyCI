using eshopPractice.ApplicationCore.Interfaces;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;

namespace eshopPractice.ApplicationCore.Entities
{
    public class CatalogBrand : BaseEntity, IAggregateRoot
    {
        public string Brand { get; private set; }
        public CatalogBrand(string brand)
        {
            Brand = brand;
        }
    }
}
