using InMemoryDbSample.Data;
using InMemoryDbSample.Entities;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;

namespace InMemoryDbSample.Services
{
    public class CarService
    {
        CarDbContext _dbContext { get; }

        public CarService(CarDbContext dbContext) => _dbContext = dbContext;

        public void Add(Car car)
        {
            _dbContext.Cars.Add(car);
            _dbContext.SaveChanges();
        }

        public IEnumerable<Car> GetById(int id) =>
            _dbContext.Cars
                .Where(car => car.Id == id)
                .OrderBy(car => car.Id)
                .ToArray();
    }
}
