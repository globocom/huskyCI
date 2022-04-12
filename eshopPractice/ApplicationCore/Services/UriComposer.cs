using eshopPractice.ApplicationCore.Interfaces;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;

namespace eshopPractice.ApplicationCore.Services
{
    public class UriComposer : IUriComposer
    {
        private readonly CatalogSettings _catalogSettings;
        public UriComposer(CatalogSettings catalogSettings) => _catalogSettings = catalogSettings;
        public string ComposePicUri(string uriTemplate)
        {
            return uriTemplate.Replace("http://catalogbaseurltobereplaced", _catalogSettings.CatalogBaseUrl);
        }
    }
}
