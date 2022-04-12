using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;

namespace eshopPractice.ApplicationCore.Entities
{
    public abstract class BaseEntity
    {
        public virtual int Id { get; protected set; }
    }
}
