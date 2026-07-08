require_relative 'titan_client'

puts "=== Cliente Ruby - Teste do 3SUM ==="
client = TitanClient.new('127.0.0.1', 8001)

begin
  client.put("a", "1")
  client.put("b", "-2")
  client.put("c", "3")
  client.put("d", "-1")
  client.put("e", "0")
  puts "Valores inseridos."

  res = client.three_sum
  if res[:status] == 0
    puts "Triplas encontradas: #{res[:data]}"
  else
    puts "Erro ao executar 3SUM (status #{res[:status]})"
  end
rescue StandardError => e
  puts "Erro: #{e.message}"
ensure
  client.close
end