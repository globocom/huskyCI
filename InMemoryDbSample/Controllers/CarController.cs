using InMemoryDbSample.Data;
using InMemoryDbSample.Entities;
using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;

namespace InMemoryDbSample.Controllers
{
    [Route("api/[controller]")]
    public class CarController : Controller
    {
        private readonly CarDbContext _dbContext;
        public CarController(CarDbContext dbContext)
        {
            _dbContext = dbContext;
        }

        [HttpGet]
        public async Task<ActionResult<IEnumerable<Car>>> Get() =>
            await _dbContext.Cars.ToListAsync();
    }
}
