class ApplicationController < ActionController::Base

    def test(options)
        system("ls #{options}") # #nohusky
    end
end

