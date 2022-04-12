using eshopPractice.ApplicationCore.Interfaces;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;

namespace eshopPractice.ApplicationCore.Entities
{
    public class CatalogType : BaseEntity, IAggregateRoot
    {
        public string Type { get; private set; }
        public CatalogType(string type)
        {
            Type = type;
        }
    }
}
