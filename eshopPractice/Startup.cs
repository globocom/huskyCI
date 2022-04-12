using eshopPractice.ApplicationCore;
using eshopPractice.ApplicationCore.Interfaces;
using eshopPractice.ApplicationCore.Services;
using eshopPractice.Infrastructure.Data;
using eshopPractice.Interfaces;
using eshopPractice.Services;
using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Hosting;
using Microsoft.AspNetCore.Http;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Microsoft.EntityFrameworkCore;
using eshopPractice.Infrastructure.Logging;
using eshopPractice.Infrastructure.Identity;
using Microsoft.AspNetCore.Identity;


namespace eshopPractice
{
    public class Startup
    {
        // This method gets called by the runtime. Use this method to add services to the container.
        // For more information on how to configure your application, visit https://go.microsoft.com/fwlink/?LinkID=398940

        public IConfiguration Configuration { get; }
        public Startup(IConfiguration configuration)
        {
            Configuration = configuration;
        }

        public void ConfigureDevelopmentServices(IServiceCollection services)
        {
            // use in-memory database
            ConfigureInMemoryDatabases(services);

            // use real database
            //ConfigureProductionServices(services);
        }

        private void ConfigureInMemoryDatabases(IServiceCollection services)
        {
            services.AddDbContext<CatalogContext>(c =>
                c.UseInMemoryDatabase("Catalog"));

            services.AddDbContext<AppIdentityDbContext>(options =>
                options.UseInMemoryDatabase("Identity"));

            ConfigureServices(services);
        }

        public void ConfigureServices(IServiceCollection services)
        {
            services.AddIdentity<ApplicationUser, IdentityRole>()
                       .AddDefaultUI()
                       .AddEntityFrameworkStores<AppIdentityDbContext>()
                       .AddDefaultTokenProviders();

            services.AddRazorPages();

            services.AddScoped(typeof(IAsyncRepository<>), typeof(EfRepository<>));
            services.AddSingleton<IUriComposer>(new UriComposer(Configuration.Get<CatalogSettings>()));
            services.AddScoped<CatalogViewModelService>();
            services.AddScoped<ICatalogItemViewModelService, CatalogItemViewModelService>();
            services.AddScoped<ICatalogViewModelService, CachedCatalogVIewModelService>();
            services.AddScoped<IBasketViewModelService, BasketViewModelService>();
            services.AddScoped<IBasketService, BasketService>();
            services.AddScoped(typeof(IAppLogger<>), typeof(LoggerAdapter<>));

            services.AddHttpContextAccessor();
        }

        // This method gets called by the runtime. Use this method to configure the HTTP request pipeline.
        public void Configure(IApplicationBuilder app, IWebHostEnvironment env)
        {
            if (env.IsDevelopment())
            {
                app.UseDeveloperExceptionPage();
            }

            app.UseStaticFiles();

            app.UseRouting();

            app.UseAuthentication();
            app.UseAuthorization();

            app.UseEndpoints(endpoints =>
            {
                endpoints.MapRazorPages();
            });
        }
    }
}
