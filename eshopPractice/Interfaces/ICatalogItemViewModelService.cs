using eshopPractice.ViewModels;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;

namespace eshopPractice.Interfaces
{
    public interface ICatalogItemViewModelService
    {
        Task UpdateCatalogItem(CatalogItemViewModel viewModel);
    }
}
