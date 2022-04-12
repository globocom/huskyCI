using InMemoryDbSample.Data;
using InMemoryDbSample.Services;
using Microsoft.EntityFrameworkCore;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Xunit;

namespace InMemoryDbSample.Test
{
    public class CarDbContextTest
    {
        public int AddExample(int x, int y)
        {
            return x + y;
        }

        [Fact]
        public void Test1()
        {
            Assert.Equal(2, AddExample(1, 1));
        }

        //[Fact]
        //public void Add_writes_to_database()
        //{
        //    var options = new DbContextOptionsBuilder<CarDbContext>()
        //        .UseInMemoryDatabase(databaseName: "Add_writes_to_database")
        //        .Options;

        //    using (var context = new CarDbContext(options))
        //    {
        //        var service = new CarService(context);
        //        service.Add(new Entities.Car
        //        {
        //            Id = 1,
        //            Name = "Accord",
        //            Descritpion = "aaaaaa"
        //        });
        //    }

        //    using (var context  =new CarDbContext(options))
        //    {
        //        Assert.Equal(1, context.Cars.Count());
        //        var car = context.Cars.Single(e => e.Id == 1);
        //        Assert.Equal("Accord", car.Name);
        //    }
        //}
    }
}
