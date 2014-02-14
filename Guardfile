guard 'shell' do
  watch(/.*\.go$/) do |m|
    puts "Starting Tests"
    `go test`
  end
end
