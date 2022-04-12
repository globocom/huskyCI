using eshopPractice.ApplicationCore.Entities;
using eshopPractice.ApplicationCore.Interfaces;
using eshopPractice.Interfaces;
using eshopPractice.ViewModels;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;

namespace eshopPractice.Services
{
    public class CatalogItemViewModelService : ICatalogItemViewModelService
    {
        private readonly IAsyncRepository<CatalogItem> _catalogItemRepository;
        
        public CatalogItemViewModelService(IAsyncRepository<CatalogItem> catalogItemRepository)
        {
            _catalogItemRepository = catalogItemRepository;
        }

        public async Task UpdateCatalogItem(CatalogItemViewModel viewModel)
        {
            var existingCatalogItem = await _catalogItemRepository.GetByIdAsync(viewModel.Id);
            existingCatalogItem.UpdateDetails(viewModel.Name, existingCatalogItem.Description, viewModel.Price);
            await _catalogItemRepository.UpdateAsync(existingCatalogItem);
        }
    }
}
