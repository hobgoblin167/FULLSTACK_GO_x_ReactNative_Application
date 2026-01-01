import React, { useState } from "react";
import {
  View,
  Text,
  TextInput,
  TouchableOpacity,
  StyleSheet,
} from "react-native";

const API_URL = "ваш_айпи:8080";

export default function RegisterScreen({ navigation }) {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [message, setMessage] = useState("");
  const [loading, setLoading] = useState(false);

  const handleRegister = async () => {
    if (!username || !password || !confirmPassword) {
      setMessage("Заполните все поля");
      return;
    }

    if (password !== confirmPassword) {
      setMessage("Пароли не совпадают");
      return;
    }

    setLoading(true);
    setMessage("");

    try {
      const response = await fetch(`${API_URL}/api/register`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          username,
          password,
        }),
      });

      const data = await response.json();
      setMessage(data.message);

      if (data.success) {
        console.log("Регистрация успешна");
      }

    } catch (error) {
      setMessage("Ошибка подключения к серверу");
    } finally {
      setLoading(false);
    }
  };

  return (
    <View style={styles.container}>
      <Text style={styles.title}>Регистрация</Text>

      <TextInput
        style={styles.input}
        placeholder="Имя пользователя"
        value={username}
        onChangeText={setUsername}
        autoCapitalize="none"
      />

      <TextInput
        style={styles.input}
        placeholder="Пароль"
        value={password}
        onChangeText={setPassword}
        secureTextEntry
      />

      <TextInput
        style={styles.input}
        placeholder="Повторите пароль"
        value={confirmPassword}
        onChangeText={setConfirmPassword}
        secureTextEntry
      />

      {message ? <Text style={styles.message}>{message}</Text> : null}

      <TouchableOpacity
        style={styles.button}
        onPress={handleRegister}
        disabled={loading}
      >
        <Text style={styles.buttonText}>
          {loading ? "Регистрация..." : "Зарегистрироваться"}
        </Text>
      </TouchableOpacity>

      <TouchableOpacity onPress={() => navigation.navigate("Login")}>
        <Text style={styles.link}>Уже есть аккаунт? Войти</Text>
      </TouchableOpacity>
    </View>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, justifyContent: "center", padding: 24 },
  title: { fontSize: 28, fontWeight: "bold", textAlign: "center", marginBottom: 24 },
  input: { borderWidth: 1, borderColor: "#ccc", borderRadius: 8, padding: 12, marginBottom: 12 },
  button: { backgroundColor: "#16a34a", padding: 14, borderRadius: 8 },
  buttonText: { color: "#fff", textAlign: "center", fontSize: 16 },
  link: { marginTop: 16, textAlign: "center", color: "#4f46e5" },
  message: { textAlign: "center", marginBottom: 12, color: "#dc2626" },
});
