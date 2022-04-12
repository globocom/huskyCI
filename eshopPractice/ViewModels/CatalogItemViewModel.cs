using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;

namespace eshopPractice.ViewModels
{
    public class CatalogItemViewModel
    {
        public int Id { get; set; }
        public string Name { get; set; }
        public string PictureUri { get; set; }
        public decimal Price { get; set; }
    }
}
