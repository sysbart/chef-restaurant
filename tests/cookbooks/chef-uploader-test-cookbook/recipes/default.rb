template "/tmp/chef-uploader" do
	source "chef-uploader.erb"
	action :create
    variables({
        :test_name => node["chef-uploader-test-cookbook"]["name"]
    })
end
