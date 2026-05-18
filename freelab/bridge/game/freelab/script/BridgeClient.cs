using Godot;
using System;
using System.Threading.Tasks;

public partial class BridgeClient : Node
{
	// Явно указываем, что используем HttpClient из System.Net.Http, а не из Godot
	private static readonly System.Net.Http.HttpClient _client = new System.Net.Http.HttpClient();
	private const string BaseUrl = "http://localhost:8080/api";

	public override void _Ready()
	{
		GD.Print("[Godot] Запуск теста связи с FreeLab Daemon...");
		
		// Запускаем асинхронный запрос при старте сцены
		_ = TestSystemStatusAsync();
	}

	private async Task TestSystemStatusAsync()
	{
		try
		{
			GD.Print($"[Godot] Отправляем GET запрос на {BaseUrl}/system/status ...");
			
			// Также явно указываем System.Net.Http для HttpResponseMessage
			System.Net.Http.HttpResponseMessage response = await _client.GetAsync($"{BaseUrl}/system/status");
			response.EnsureSuccessStatusCode(); // Проверяем, что ответ 200 OK
			
			// Читаем JSON ответ
			string responseBody = await response.Content.ReadAsStringAsync();
			
			GD.Print("[Godot] УСПЕХ! Ответ от ядра (Go):");
			GD.Print(responseBody);
		}
		catch (System.Net.Http.HttpRequestException e)
		{
			GD.PrintErr($"[Godot] ОШИБКА СВЯЗИ! Убедитесь, что Go-сервер запущен. Детали: {e.Message}");
		}
	}
}
