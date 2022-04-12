using System;
using System.Collections.Generic;
using System.Linq;
using System.Security.Cryptography.X509Certificates;
using System.Threading.Tasks;

namespace InMemoryDbSample.Entities
{
    public class Car
    {
        public int Id { get; set; }
        public string Name { get; set; }
        public string Descritpion { get; set; }
    }
}
